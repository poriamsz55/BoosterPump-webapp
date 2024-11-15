package database

import (
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
