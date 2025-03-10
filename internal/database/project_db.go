package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
	tehrantime "github.com/poriamsz55/BoosterPump-webapp/internal/time"
)

func AddProjectToDB(p *project.Project) (int, error) {

	// check if the project exists
	err := CheckProjectFromDB(p.Name)
	if err != nil {
		return -1, errors.New("Project already exists")
	}

	query := `INSERT INTO ` + tableProjects + ` (` + columnProjectName + `) 
	          VALUES (?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.Name)
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

func CheckProjectFromDB(name string) error {
	// check if the project exists
	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ?
	`, tableProjects, columnProjectName)

	var count int
	err := instance.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("Project already exists")
	}

	return nil
}

func GetAllProjects(db *sql.DB) ([]*project.Project, error) {

	if db == nil {
		db = instance.db
	}

	query := fmt.Sprintf(`
        SELECT %s, %s, %s
        FROM %s
    `, columnProjectID, columnProjectName, columnModifiedAt,
		tableProjects)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*project.Project
	for rows.Next() {
		p := project.NewEmptyProject()
		err := rows.Scan(&p.Id, &p.Name, &p.ModifiedAt)
		if err != nil {
			return nil, err
		}

		// Get devices for this project
		p.ProjectDeviceList, err = GetProjectDevicesByProjectId(db, p.Id)
		if err != nil {
			return nil, err
		}

		// Get extra prices for this project
		p.ExtraPriceList, err = GetExtraPricesByProjectIdFromDB(db, p.Id)
		if err != nil {
			return nil, err
		}

		p.UpdatePrice()

		projects = append(projects, p)
	}

	return projects, nil
}

func GetProjectByIdFromDB(db *sql.DB, id int) (*project.Project, error) {

	if db == nil {
		db = instance.db
	}

	query := fmt.Sprintf(`
        SELECT %s, %s, %s
        FROM %s
        WHERE %s = ?
    `, columnProjectID, columnProjectName, columnModifiedAt,
		tableProjects, columnProjectID)

	p := project.NewEmptyProject()
	err := instance.db.QueryRow(query, id).Scan(&p.Id, &p.Name, &p.ModifiedAt)
	if err != nil {
		return nil, err
	}

	// Get devices for this project
	p.ProjectDeviceList, err = GetProjectDevicesByProjectId(db, id)
	if err != nil {
		return nil, err
	}

	// Get extra prices for this project
	p.ExtraPriceList, err = GetExtraPricesByProjectIdFromDB(db, id)
	if err != nil {
		return nil, err
	}

	p.UpdatePrice()

	p.ModifiedAt = tehrantime.FormattedDateTime(p.ModifiedAt)
	return p, nil
}

func DeleteProjectFromDB(id int) error {

	query := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableProjects, columnProjectID)

	result, err := instance.db.Exec(query, id)
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

	// delete related records in project_devices table
	err = DeleteProjectDevicesByProjectId(id)
	if err != nil {
		return err
	}

	// delete related records in extra_prices table
	err = DeleteExtraPricesByProjectId(id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProjectInDB(id int, name string, projectDevices []device.DeviceReq) error {

	// First check if the project exists
	checkQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE %s = ?
	`, tableProjects, columnProjectID)

	var count int
	err := instance.db.QueryRow(checkQuery, id).Scan(&count)
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

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?
        WHERE %s = ?
    `, tableProjects, columnProjectName, columnProjectID)

	result, err := instance.db.Exec(query, name, id)
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

	// Delete project devices
	deleteDevicesQuery := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableProjectDevices, columnProjectIDFK)

	_, err = tx.Exec(deleteDevicesQuery, id)
	if err != nil {
		return err
	}

	// Insert project devices
	if len(projectDevices) > 0 {
		insertDevicesQuery := fmt.Sprintf(`
		INSERT INTO %s (%s, %s, %s)
		VALUES (?, ?, ?)
	`, tableProjectDevices, columnProjectDeviceCount, columnDeviceIDFK, columnProjectIDFK)

		for _, dp := range projectDevices {
			_, err = tx.Exec(insertDevicesQuery,
				dp.Count,
				dp.Id,
				id)
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
