package database

import "log"

func AddDevicePartToDB(dvcId int, count float32, prtId int) error {
	query := `INSERT INTO ` + tableDeviceParts + ` (` + columnDevicePartCount + `, ` +
		columnPartIDK + `, ` +
		columnDeviceIDFK + `) 
	          VALUES (?, ?, ?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dvcId, count, prtId)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}
