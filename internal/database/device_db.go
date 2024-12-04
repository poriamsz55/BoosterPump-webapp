package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
)

func AddDeviceToDB(d *device.Device) (int, error) {
	query := `INSERT INTO ` + tableDevices + ` (` +
		columnDeviceName + `, ` +
		columnDeviceConverter + `, ` +
		columnDeviceFilter + `) 
	          VALUES (?, ?, ?)`

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

	result, err := stmt.Exec(d.Name, d.Converter.String(), filterInt)
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

func GetAllDevicesFromDB() ([]*device.Device, error) {
	// First, get all devices
	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s 
        FROM %s
    `, columnDeviceID, columnDeviceName, columnDeviceConverter, columnDeviceFilter, tableDevices)

	rows, err := instance.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*device.Device
	for rows.Next() {
		var id int
		var name string
		var converter string
		var filter bool

		err := rows.Scan(&id, &name, &converter, &filter)
		if err != nil {
			return nil, err
		}

		converterConv, err := device.ConverterFromName(converter)
		if err != nil {
			return nil, err
		}

		dev := device.NewDevice(name, converterConv, filter)
		dev.Id = id

		// Get device parts for this device
		deviceParts, err := GetDevicePartsByDeviceId(id)
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

func GetDeviceByIdFromDB(deviceID int) (*device.Device, error) {

	// Get device information
	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s 
        FROM %s 
        WHERE %s = ?
    `, columnDeviceID, columnDeviceName, columnDeviceConverter, columnDeviceFilter,
		tableDevices, columnDeviceID)

	row := instance.db.QueryRow(query, deviceID)

	var id int
	var name string
	var converterStr string
	var filter bool

	err := row.Scan(&id, &name, &converterStr, &filter)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	converter, err := device.ConverterFromName(converterStr)
	if err != nil {
		return nil, err
	}

	dev := device.NewDevice(name, converter, filter)
	dev.Id = id

	// Get device parts for this device
	deviceParts, err := GetDevicePartsByDeviceId(id)
	if err != nil {
		return nil, err
	}
	dev.DevicePartList = deviceParts
	dev.UpdatePrice()

	return dev, nil
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

	return nil
}

func UpdateDeviceInDB(updatedDevice *device.Device) error {

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

	converter, err := device.ConverterFromValue(int(updatedDevice.Converter))
	if err != nil {
		return err
	}

	var filterInt int
	if updatedDevice.Filter {
		filterInt = 1
	} else {
		filterInt = 0
	}

	result, err := tx.Exec(query,
		updatedDevice.Name,
		converter,
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
	if len(updatedDevice.DevicePartList) > 0 {
		insertPartsQuery := fmt.Sprintf(`
            INSERT INTO %s (%s, %s, %s) 
            VALUES (?, ?, ?)
        `, tableDeviceParts,
			columnDevicePartCount, columnPartIDK, columnDeviceIDFK)

		for _, dp := range updatedDevice.DevicePartList {
			_, err = tx.Exec(insertPartsQuery,
				dp.Count,
				dp.Part.Id,
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
