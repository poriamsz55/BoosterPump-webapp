package loader

import (
	"database/sql"
	"fmt"
)

// Struct for Project
type Project struct {
	ID   int64
	Name string
}

func convertProjects(javaDB, goDB *sql.DB) error {

	// Query all projects from Java DB
	rows, err := javaDB.Query(`
        SELECT project_id, project_name 
        FROM projects
    `)
	if err != nil {
		return fmt.Errorf("error querying Java database: %v", err)
	}
	defer rows.Close()

	// Begin transaction for Go DB
	tx, err := goDB.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback() // Will rollback if not committed

	// Prepare the insert statement
	stmt, err := tx.Prepare(`
        INSERT INTO projects (project_id, project_name) 
        VALUES (?, ?)
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Process each project
	var count int
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
		)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Insert into Go database
		_, err = stmt.Exec(
			project.ID,
			project.Name,
		)
		if err != nil {
			return fmt.Errorf("error inserting project %d: %v", project.ID, err)
		}

		count++
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	fmt.Printf("Converted %d projects\n", count)
	return nil
}

// Utility function to verify the conversion
func verifyProjectsConversion(javaDB, goDB *sql.DB) error {
	var javaCount, goCount int

	// Count records in Java DB
	err := javaDB.QueryRow("SELECT COUNT(*) FROM projects").Scan(&javaCount)
	if err != nil {
		return fmt.Errorf("error counting Java records: %v", err)
	}

	// Count records in Go DB
	err = goDB.QueryRow("SELECT COUNT(*) FROM projects").Scan(&goCount)
	if err != nil {
		return fmt.Errorf("error counting Go records: %v", err)
	}

	fmt.Printf("Projects count - Java DB: %d, Go DB: %d\n", javaCount, goCount)

	if javaCount != goCount {
		return fmt.Errorf("record count mismatch: Java=%d, Go=%d", javaCount, goCount)
	}

	return nil
}
