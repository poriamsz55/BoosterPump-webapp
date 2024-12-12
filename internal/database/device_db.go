package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
	tehrantime "github.com/poriamsz55/BoosterPump-webapp/internal/time"
)

func AddDeviceToDB(d *device.Device) (int, error) {
	query := `INSERT INTO ` + tableDevices + ` (` +
		columnDeviceName + `, ` +
		columnDeviceConverter + `, ` +
		columnDeviceFilter + `) 
	          VALUES (?, ?, ?)`

	// check if device exists
	err := CheckDeviceFromDB(d)
	if err != nil {
		return -1, err
	}

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return -1, err
	}
	defer stmt.Close()

	var filterInt int
	if d.Filter {
		filterInt = 1
	} else {
		filterInt = 0
	}

	result, err := stmt.Exec(d.Name, int(d.Converter), filterInt)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert id: %v", err)
		return -1, err
	}

	return int(id), nil
}

func CheckDeviceFromDB(d *device.Device) error {
	// First check if the device exists
	// check by name and converter and filter
	query := `
		SELECT COUNT(*) 
		FROM ` + tableDevices + `
		WHERE ` + columnDeviceName + ` = ? AND ` + columnDeviceConverter + ` = ? AND ` + columnDeviceFilter + ` = ?
	`

	var count int
	err := instance.db.QueryRow(query, d.Name, int(d.Converter), d.Filter).Scan(&count)
	if err != nil {
		log.Printf("Error checking device: %v", err)
		return err
	}
	if count > 0 {
		log.Printf("Device already exists")
		return errors.New("Device already exists")
	}

	return nil
}

func GetAllDevices(db *sql.DB) ([]*device.Device, error) {

	if db == nil {
		db = instance.db
	}

	// First, get all devices
	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s, %s
        FROM %s
    `, columnDeviceID, columnDeviceName, columnDeviceConverter, columnDeviceFilter, columnModifiedAt,
		tableDevices)

	rows, err := instance.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*device.Device
	for rows.Next() {
		var id int
		var name string
		var converter int
		var filter bool
		var mt time.Time

		err := rows.Scan(&id, &name, &converter, &filter, &mt)
		if err != nil {
			return nil, err
		}

		converterConv, err := device.ConverterFromValue(converter)
		if err != nil {
			return nil, err
		}

		dev := device.NewDevice(name, converterConv, filter)
		dev.Id = id
		dev.ModifiedAt = tehrantime.FormattedDateTime(mt)

		// Get device parts for this device
		deviceParts, err := GetDevicePartsByDeviceId(db, id)
		if err != nil {
			return nil, err
		}
		dev.DevicePartList = deviceParts
		dev.UpdatePrice()

		devices = append(devices, dev)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func GetDeviceByIdFromDB(db *sql.DB, deviceID int) (*device.Device, error) {

	if db == nil {
		db = instance.db
	}

	// Get device information
	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s, %s
        FROM %s 
        WHERE %s = ?
    `, columnDeviceID, columnDeviceName, columnDeviceConverter, columnDeviceFilter, columnModifiedAt,
		tableDevices, columnDeviceID)

	row := db.QueryRow(query, deviceID)

	var id int
	var name string
	var converterInt int
	var filter bool
	var mt time.Time

	err := row.Scan(&id, &name, &converterInt, &filter, &mt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	converter, err := device.ConverterFromValue(converterInt)
	if err != nil {
		return nil, err
	}

	dev := device.NewDevice(name, converter, filter)
	dev.Id = id
	dev.ModifiedAt = tehrantime.FormattedDateTime(mt)

	// Get device parts for this device
	deviceParts, err := GetDevicePartsByDeviceId(db, id)
	if err != nil {
		return nil, err
	}
	dev.DevicePartList = deviceParts
	dev.UpdatePrice()

	return dev, nil
}

// Check if Device exists By Name
func CheckDeviceByNameFromDB(name string) error {

	// First check if the device exists
	checkQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ?
	`, tableDevices, columnDeviceName)

	var count int
	err := instance.db.QueryRow(checkQuery, name).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteDeviceFromDB(deviceID int) error {

	// First check if the device exists
	checkQuery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s 
        WHERE %s = ?
    `, tableDevices, columnDeviceID)

	var count int
	err := instance.db.QueryRow(checkQuery, deviceID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	// Since we have ON DELETE CASCADE in our foreign keys,
	// deleting from the devices table will automatically delete
	// related records in device_parts and project_devices tables
	query := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableDevices, columnDeviceID)

	result, err := instance.db.Exec(query, deviceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Delete device parts for this device
	deleteDevicePartsQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = ?
	`, tableDeviceParts, columnDeviceID)
	_, err = instance.db.Exec(deleteDevicePartsQuery, deviceID)
	if err != nil {
		return err
	}

	// Delete project devices for this device
	deleteProjectDevicesQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = ?
	`, tableProjectDevices, columnDeviceID)
	_, err = instance.db.Exec(deleteProjectDevicesQuery, deviceID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDeviceInDB(updatedDevice *device.Device, partsReq []part.PartReq) error {

	// First check if the device exists
	checkQuery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s 
        WHERE %s = ?
    `, tableDevices, columnDeviceID)

	var count int
	err := instance.db.QueryRow(checkQuery, updatedDevice.Id).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	// Begin transaction
	tx, err := instance.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Update device information
	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?, %s = ?, %s = ?
        WHERE %s = ?
    `, tableDevices,
		columnDeviceName, columnDeviceConverter, columnDeviceFilter,
		columnDeviceID)

	var filterInt int
	if updatedDevice.Filter {
		filterInt = 1
	} else {
		filterInt = 0
	}

	result, err := tx.Exec(query,
		updatedDevice.Name,
		int(updatedDevice.Converter),
		filterInt,
		updatedDevice.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Delete existing device parts
	deletePartsQuery := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableDeviceParts, columnDeviceIDFK)

	_, err = tx.Exec(deletePartsQuery, updatedDevice.Id)
	if err != nil {
		return err
	}

	// Insert new device parts if any
	if len(partsReq) > 0 {
		insertPartsQuery := fmt.Sprintf(`
            INSERT INTO %s (%s, %s, %s) 
            VALUES (?, ?, ?)
        `, tableDeviceParts,
			columnDevicePartCount, columnPartIDK, columnDeviceIDFK)

		for _, dp := range partsReq {
			_, err = tx.Exec(insertPartsQuery,
				dp.Count,
				dp.Id,
				updatedDevice.Id)
			if err != nil {
				return err
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
