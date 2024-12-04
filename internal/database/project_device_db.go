package database

import (
	"fmt"
	"log"

	projectd "github.com/poriamsz55/BoosterPump-webapp/internal/models/project_device"
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

func GetProjectDevicesByProjectId(projectID int) ([]*projectd.ProjectDevice, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, tableProjectDevices, columnProjectIDFK)

	rows, err := instance.db.Query(query, projectID)
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
		)
		if err != nil {
			return nil, err
		}

		// get device by device Id
		dvc, err := GetDeviceByIdFromDB(dvcId)
		if err != nil {
			return nil, err
		}
		p.Device = dvc

		projectDevices = append(projectDevices, &p)
	}

	return projectDevices, nil
}
