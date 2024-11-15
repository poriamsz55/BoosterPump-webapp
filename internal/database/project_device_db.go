package database

import (
	"log"
)

func AddProjectDeviceToDB(prjId int, count float32, dvcId int) error {
	query := `INSERT INTO ` + tableProjectDevices + ` (` + columnProjectDeviceCount + `, ` +
		columnDeviceIDK + `, ` +
		columnProjectIDFK + `) 
	          VALUES (?, ?, ?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(prjId, count, dvcId)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}
