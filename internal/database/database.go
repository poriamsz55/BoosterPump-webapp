package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver

	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	driveName    = "sqlite3"
	databaseName = "booster_pump.db"

	tableProjects       = "projects"
	tableExtraPrice     = "extra_price"
	tableProjectDevices = "project_devices"
	tableDeviceParts    = "device_parts"
	tableDevices        = "devices"
	tableParts          = "parts"

	columnProjectID          = "project_id"
	columnProjectName        = "project_name"
	columnExtraPriceID       = "extra_price_id"
	columnExtraPriceName     = "extra_price_name"
	columnExtraPriceValue    = "extra_price_value"
	columnProjectDeviceID    = "project_device_id"
	columnProjectDeviceCount = "project_device_count"
	columnProjectIDFK        = "project_id"
	columnDeviceIDK          = "device_id"
	columnDevicePartID       = "device_part_id"
	columnDevicePartCount    = "device_part_count"
	columnDeviceIDFK         = "device_id"
	columnPartIDK            = "part_id"
	columnDeviceID           = "device_id"
	columnDeviceName         = "device_name"
	columnDeviceConverter    = "device_converter"
	columnDeviceFilter       = "device_filter"
	columnPartID             = "part_id"
	columnPartName           = "part_name"
	columnPartSize           = "part_size"
	columnPartMaterial       = "part_material"
	columnPartBrand          = "part_brand"
	columnPartPrice          = "part_price"
)

var (
	createProjectsTable = fmt.Sprintf(`CREATE TABLE %s (
        %s INTEGER PRIMARY KEY AUTOINCREMENT,
        %s TEXT
    )`, tableProjects, columnProjectID, columnProjectName)

	createExtraPriceTable = fmt.Sprintf(`CREATE TABLE %s (
        %s INTEGER PRIMARY KEY AUTOINCREMENT,
        %s TEXT,
        %s INTEGER,
        %s INTEGER,
        FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE CASCADE
    )`, tableExtraPrice, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue, columnProjectIDFK, columnProjectIDFK, tableProjects, columnProjectID)

	createProjectDevicesTable = fmt.Sprintf(`CREATE TABLE %s (
        %s INTEGER PRIMARY KEY AUTOINCREMENT,
        %s FLOAT,
        %s INTEGER,
        %s INTEGER,
        FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE CASCADE
    )`, tableProjectDevices, columnProjectDeviceID, columnProjectDeviceCount, columnDeviceIDK, columnProjectIDFK, columnDeviceIDK, tableDevices, columnDeviceID, columnProjectIDFK, tableProjects, columnProjectID)

	createDevicePartsTable = fmt.Sprintf(`CREATE TABLE %s (
        %s INTEGER PRIMARY KEY AUTOINCREMENT,
        %s FLOAT,
        %s INTEGER,
        %s INTEGER,
        FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY(%s) REFERENCES %s(%s) ON DELETE CASCADE ON UPDATE CASCADE
    )`, tableDeviceParts, columnDevicePartID, columnDevicePartCount, columnPartIDK, columnDeviceIDFK, columnPartIDK, tableParts, columnPartID, columnDeviceIDFK, tableDevices, columnDeviceID)

	createDevicesTable = fmt.Sprintf(`CREATE TABLE %s (
        %s INTEGER PRIMARY KEY AUTOINCREMENT,
        %s TEXT,
        %s INTEGER,
        %s INTEGER
    )`, tableDevices, columnDeviceID, columnDeviceName, columnDeviceConverter, columnDeviceFilter)

	createPartsTable = fmt.Sprintf(`CREATE TABLE %s (
        %s INTEGER PRIMARY KEY AUTOINCREMENT,
        %s TEXT,
        %s TEXT,
        %s TEXT,
        %s TEXT,
        %s INTEGER
    )`, tableParts, columnPartID, columnPartName, columnPartSize, columnPartMaterial, columnPartBrand, columnPartPrice)
)

type DBHelper struct {
	db           *sql.DB
	databasePath string
	mu           sync.Mutex
}

var instance *DBHelper

func GetDBHelperInstance() *DBHelper {
	if instance != nil {
		return instance
	}

	instance = &DBHelper{
		databasePath: filepath.Join("./", databaseName),
	}
	return instance
}

func (h *DBHelper) Open() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var err error
	h.db, err = sql.Open(driveName, h.databasePath)
	if err != nil {
		return err
	}

	_, err = h.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	return nil
}

func (h *DBHelper) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

func (h *DBHelper) CheckDatabase() bool {
	_, err := os.Stat(h.databasePath)
	return err == nil
}

func (h *DBHelper) CreateTables() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.db.Exec(createProjectsTable)
	if err != nil {
		return err
	}

	_, err = h.db.Exec(createExtraPriceTable)
	if err != nil {
		return err
	}

	_, err = h.db.Exec(createProjectDevicesTable)
	if err != nil {
		return err
	}

	_, err = h.db.Exec(createDevicePartsTable)
	if err != nil {
		return err
	}

	_, err = h.db.Exec(createDevicesTable)
	if err != nil {
		return err
	}

	_, err = h.db.Exec(createPartsTable)
	if err != nil {
		return err
	}

	return nil
}

func (h *DBHelper) DropTables() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	tables := []string{tableParts, tableDevices, tableDeviceParts, tableProjectDevices, tableExtraPrice, tableProjects}
	for _, table := range tables {
		_, err := h.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			return err
		}
	}

	return nil
}

func CloseDB() {
	dbHelper := GetDBHelperInstance()
	dbHelper.Close()
}

func InitializeDB() {
	dbHelper := GetDBHelperInstance()

	// check if database exists
	databaseExists := dbHelper.CheckDatabase()

	err := dbHelper.Open()
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	err = dbHelper.CreateTables()
	if err != nil {
		fmt.Println("Error creating tables:", err)
		return
	}

	if !databaseExists {
		err = dbHelper.LoadInitDatabase()
		if err != nil {
			fmt.Println("Error copying database from assets:", err)
			return
		}
	}

}
