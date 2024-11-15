package database

import (
	"log"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
)

func AddDeviceToDB(d *device.Device) error {
	query := `INSERT INTO ` + tableDevices + ` (` +
		columnDeviceName + `, ` +
		columnDeviceConverter + `, ` +
		columnDeviceFilter + `) 
	          VALUES (?, ?, ?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	var filterInt int
	if d.Filter {
		filterInt = 1
	} else {
		filterInt = 0
	}

	_, err = stmt.Exec(d.Name, d.Converter.String(), filterInt)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}
