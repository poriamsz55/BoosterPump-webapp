package loader

import (
	"database/sql"
	"fmt"
	"log"
)

// Struct to represent a device from Java DB
type JavaDevice struct {
	ID        int64
	Name      string
	Converter string // In Java it's TEXT
	Filter    int
}

// Struct to represent a device for Go DB
type GoDevice struct {
	ID        int64
	Name      string
	Converter int // In Go it's INTEGER
	Filter    int
}

func main() {
	// Open the Java SQLite database (init.db)
	javaDB, err := sql.Open("sqlite3", "./init.db")
	if err != nil {
		log.Fatal("Error opening Java database:", err)
	}
	defer javaDB.Close()

	// Open the Go SQLite database
	goDB, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal("Error opening Go database:", err)
	}
	defer goDB.Close()

	// Convert the data
	if err := convertDevices(javaDB, goDB); err != nil {
		log.Fatal("Error converting devices:", err)
	}

	fmt.Println("Conversion completed successfully!")
}

func convertDevices(javaDB, goDB *sql.DB) error {
	// First, read all devices from Java DB
	rows, err := javaDB.Query(`SELECT device_id, device_name, device_converter, device_filter FROM devices`)
	if err != nil {
		return fmt.Errorf("error querying Java database: %v", err)
	}
	defer rows.Close()

	// Begin transaction for Go DB
	tx, err := goDB.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback() // Will rollback if not committed

	// Prepare the insert statement
	stmt, err := tx.Prepare(`INSERT INTO devices (device_id, device_name, device_converter, device_filter) 
                            VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Process each device
	var devices []JavaDevice
	for rows.Next() {
		var dev JavaDevice
		err := rows.Scan(&dev.ID, &dev.Name, &dev.Converter, &dev.Filter)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}
		devices = append(devices, dev)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %v", err)
	}

	// Insert devices into Go database
	for _, jDev := range devices {
		gDev := convertJavaToGoDevice(jDev)

		_, err := stmt.Exec(gDev.ID, gDev.Name, gDev.Converter, gDev.Filter)
		if err != nil {
			return fmt.Errorf("error inserting device %d: %v", gDev.ID, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	// Print conversion statistics
	fmt.Printf("Converted %d devices\n", len(devices))
	return nil
}

func convertJavaToGoDevice(jDev JavaDevice) GoDevice {
	// Convert the Converter field from TEXT to INTEGER
	var converterInt int
	switch jDev.Converter {
	case "بدون تبدیل":
		converterInt = 0
	case "تبدیل در دهش":
		converterInt = 1
	case "تبدیل دو طرفه":
		converterInt = 2
	default:
		converterInt = -1 // Default value
		fmt.Printf("Warning: Unknown converter value: %s for device ID: %d\n",
			jDev.Converter, jDev.ID)
	}

	return GoDevice{
		ID:        jDev.ID,
		Name:      jDev.Name,
		Converter: converterInt,
		Filter:    jDev.Filter,
	}
}

// Utility function to verify the conversion
func verifyDevicesConversion(javaDB, goDB *sql.DB) error {
	var javaCount, goCount int

	// Count records in Java DB
	err := javaDB.QueryRow("SELECT COUNT(*) FROM devices").Scan(&javaCount)
	if err != nil {
		return fmt.Errorf("error counting Java records: %v", err)
	}

	// Count records in Go DB
	err = goDB.QueryRow("SELECT COUNT(*) FROM devices").Scan(&goCount)
	if err != nil {
		return fmt.Errorf("error counting Go records: %v", err)
	}

	fmt.Printf("Records count - Java DB: %d, Go DB: %d\n", javaCount, goCount)

	if javaCount != goCount {
		return fmt.Errorf("record count mismatch: Java=%d, Go=%d", javaCount, goCount)
	}

	return nil
}
