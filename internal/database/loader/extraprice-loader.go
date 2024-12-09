package loader

import (
	"database/sql"
	"fmt"
)

// Struct for Java ExtraPrice (source)
type JavaExtraPrice struct {
	ID        int64
	Name      string
	Value     float64 // FLOAT in Java
	ProjectID int64
}

// Struct for Go ExtraPrice (destination)
type GoExtraPrice struct {
	ID        int64
	Name      string
	Value     int // INTEGER in Go
	ProjectID int64
}

func convertExtraPrices(javaDB, goDB *sql.DB) error {
	// Query all extra prices from Java DB
	rows, err := javaDB.Query(`
        SELECT extra_price_id, extra_price_name, extra_price_value, project_id 
        FROM extra_price
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
        INSERT INTO extra_price (
            extra_price_id, 
            extra_price_name, 
            extra_price_value, 
            project_id
        ) VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Process each extra price
	var count int
	for rows.Next() {
		var jExtraPrice JavaExtraPrice
		err := rows.Scan(
			&jExtraPrice.ID,
			&jExtraPrice.Name,
			&jExtraPrice.Value,
			&jExtraPrice.ProjectID,
		)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Convert Java extra price to Go extra price
		gExtraPrice := convertJavaToGoExtraPrice(jExtraPrice)

		// Insert into Go database
		_, err = stmt.Exec(
			gExtraPrice.ID,
			gExtraPrice.Name,
			gExtraPrice.Value,
			gExtraPrice.ProjectID,
		)
		if err != nil {
			return fmt.Errorf("error inserting extra price %d: %v", gExtraPrice.ID, err)
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

	fmt.Printf("Converted %d extra prices\n", count)
	return nil
}

func convertJavaToGoExtraPrice(jExtraPrice JavaExtraPrice) GoExtraPrice {
	// Convert float value to integer (multiply by 100 to preserve 2 decimal places)
	valueInt := int(jExtraPrice.Value * 100)

	return GoExtraPrice{
		ID:        jExtraPrice.ID,
		Name:      jExtraPrice.Name,
		Value:     valueInt,
		ProjectID: jExtraPrice.ProjectID,
	}
}

func verifyExtraPricesConversion(javaDB, goDB *sql.DB) error {
	var javaCount, goCount int

	// Count records in Java DB
	err := javaDB.QueryRow("SELECT COUNT(*) FROM extra_price").Scan(&javaCount)
	if err != nil {
		return fmt.Errorf("error counting Java records: %v", err)
	}

	// Count records in Go DB
	err = goDB.QueryRow("SELECT COUNT(*) FROM extra_price").Scan(&goCount)
	if err != nil {
		return fmt.Errorf("error counting Go records: %v", err)
	}

	fmt.Printf("Extra prices count - Java DB: %d, Go DB: %d\n", javaCount, goCount)

	if javaCount != goCount {
		return fmt.Errorf("record count mismatch: Java=%d, Go=%d", javaCount, goCount)
	}

	// Verify some sample records
	rows, err := javaDB.Query(`
        SELECT ep.extra_price_id, ep.extra_price_name, ep.extra_price_value, p.project_name
        FROM extra_price ep
        JOIN projects p ON ep.project_id = p.project_id
        LIMIT 5
    `)
	if err != nil {
		return fmt.Errorf("error querying sample records: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nSample extra prices comparison:")
	for rows.Next() {
		var id int64
		var name string
		var javaValue float64
		var projectName string

		err := rows.Scan(&id, &name, &javaValue, &projectName)
		if err != nil {
			return fmt.Errorf("error scanning Java record: %v", err)
		}

		var goValue int
		err = goDB.QueryRow("SELECT extra_price_value FROM extra_price WHERE extra_price_id = ?", id).Scan(&goValue)
		if err != nil {
			return fmt.Errorf("error querying Go record: %v", err)
		}

		fmt.Printf("ID: %d, Name: %s, Project: %s\n", id, name, projectName)
		fmt.Printf("  Java value: %.2f\n", javaValue)
		fmt.Printf("  Go value: %.2f\n", float64(goValue)/100)
	}

	return nil
}
