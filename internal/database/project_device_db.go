package database

import (
	"database/sql"
	"fmt"
	"log"

	projectd "github.com/poriamsz55/BoosterPump-webapp/internal/models/project_device"
	tehrantime "github.com/poriamsz55/BoosterPump-webapp/internal/time"
)

func AddProjectDeviceToDB(prjId int, count float64, dvcId int) error {
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

	_, err = stmt.Exec(count, dvcId, prjId)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

func GetProjectDevicesByProjectId(db *sql.DB, projectID int) ([]*projectd.ProjectDevice, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, tableProjectDevices, columnProjectIDFK)

	if db == nil {
		db = instance.db
	}

	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectDevices []*projectd.ProjectDevice
	var dvcId int
	for rows.Next() {
		var p projectd.ProjectDevice
		err := rows.Scan(
			&p.Id,
			&p.Count,
			&dvcId,
			&p.ProjectId,
			&p.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		// get device by device Id
		dvc, err := GetDeviceByIdFromDB(db, dvcId)
		if err != nil {
			return nil, err
		}
		p.Device = dvc

		p.ModifiedAt = tehrantime.FormattedDateTime(p.ModifiedAt)
		projectDevices = append(projectDevices, &p)
	}

	return projectDevices, nil
}

func DeleteProjectDeviceFromDB(id int) error {
	query := `DELETE FROM ` + tableProjectDevices + ` WHERE ` + columnProjectDeviceID + ` = ?`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

func DeleteProjectDevicesByProjectId(projectID int) error {
	query := `DELETE FROM ` + tableProjectDevices + ` WHERE ` + columnProjectIDFK + ` = ?`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(projectID)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}
