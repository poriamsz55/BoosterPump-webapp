package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
)

func AddProjectToDB(p *project.Project) error {
	query := `INSERT INTO ` + tableProjects + ` (` + columnProjectName + `) 
	          VALUES (?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Name)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

func GetAllProjectsFromDB() ([]*project.Project, error) {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return nil, err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        SELECT %s, %s
        FROM %s
    `, columnProjectID, columnProjectName, tableProjects)

	rows, err := dbHelper.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*project.Project
	for rows.Next() {
		p := project.NewEmptyProject()
		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			return nil, err
		}

		// Get devices for this project
		p.ProjectDeviceList, err = GetProjectDevicesByProjectId(p.Id)
		if err != nil {
			return nil, err
		}

		// Get extra prices for this project
		p.ExtraPriceList, err = GetExtraPricesByProjectIdFromDB(p.Id)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func GetProjectByIdFromDB(id int) (*project.Project, error) {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return nil, err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        SELECT %s, %s
        FROM %s
        WHERE %s = ?
    `, columnProjectID, columnProjectName, tableProjects, columnProjectID)

	p := project.NewEmptyProject()
	err = dbHelper.db.QueryRow(query, id).Scan(&p.Id, &p.Name)
	if err != nil {
		return nil, err
	}

	// Get devices for this project
	p.ProjectDeviceList, err = GetProjectDevicesByProjectId(id)
	if err != nil {
		return nil, err
	}

	// Get extra prices for this project
	p.ExtraPriceList, err = GetExtraPricesByProjectIdFromDB(id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func DeleteProjectFromDB(id int) error {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableProjects, columnProjectID)

	result, err := dbHelper.db.Exec(query, id)
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

func UpdateProjectInDB(id int, name string) error {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?
        WHERE %s = ?
    `, tableProjects, columnProjectName, columnProjectID)

	result, err := dbHelper.db.Exec(query, name, id)
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
