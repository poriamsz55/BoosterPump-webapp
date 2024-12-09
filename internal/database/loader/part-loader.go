package loader

import (
	"database/sql"
	"fmt"
	"math"
)

// Struct for Java Part (source)
type JavaPart struct {
	ID       int
	Name     string
	Size     string
	Material string
	Brand    string
	Price    float64 // FLOAT in Java
}

// Struct for Go Part (destination)
type GoPart struct {
	ID       int
	Name     string
	Size     string
	Material string
	Brand    string
	Price    uint64 // INTEGER in Go
}

func convertParts(javaDB, goDB *sql.DB) error {
	// Query all parts from Java DB
	rows, err := javaDB.Query(`
        SELECT part_id, part_name, part_size, part_material, part_brand, part_price 
        FROM parts
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
        INSERT INTO parts (part_id, part_name, part_size, part_material, part_brand, part_price) 
        VALUES (?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Process each part
	var count int
	for rows.Next() {
		var jPart JavaPart
		err := rows.Scan(
			&jPart.ID,
			&jPart.Name,
			&jPart.Size,
			&jPart.Material,
			&jPart.Brand,
			&jPart.Price,
		)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Convert Java part to Go part
		gPart := convertJavaToGoPart(jPart)

		// Insert into Go database
		_, err = stmt.Exec(
			gPart.ID,
			gPart.Name,
			gPart.Size,
			gPart.Material,
			gPart.Brand,
			gPart.Price,
		)
		if err != nil {
			return fmt.Errorf("error inserting part %d: %v", gPart.ID, err)
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

	fmt.Printf("Converted %d parts\n", count)
	return nil
}

func convertJavaToGoPart(jPart JavaPart) GoPart {
	// Convert float price to integer (multiply by 100 to preserve 2 decimal places)
	priceInt := uint64(math.Ceil(jPart.Price))

	return GoPart{
		ID:       jPart.ID,
		Name:     jPart.Name,
		Size:     jPart.Size,
		Material: jPart.Material,
		Brand:    jPart.Brand,
		Price:    priceInt,
	}
}

func verifyPartsConversion(javaDB, goDB *sql.DB) error {
	var javaCount, goCount int

	// Count records in Java DB
	err := javaDB.QueryRow("SELECT COUNT(*) FROM parts").Scan(&javaCount)
	if err != nil {
		return fmt.Errorf("error counting Java records: %v", err)
	}

	// Count records in Go DB
	err = goDB.QueryRow("SELECT COUNT(*) FROM parts").Scan(&goCount)
	if err != nil {
		return fmt.Errorf("error counting Go records: %v", err)
	}

	fmt.Printf("Parts count - Java DB: %d, Go DB: %d\n", javaCount, goCount)

	if javaCount != goCount {
		return fmt.Errorf("record count mismatch: Java=%d, Go=%d", javaCount, goCount)
	}

	// Optional: Verify some sample records
	rows, err := javaDB.Query(`
        SELECT part_id, part_name, part_price 
        FROM parts 
        LIMIT 5
    `)
	if err != nil {
		return fmt.Errorf("error querying sample records: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nSample records comparison:")
	for rows.Next() {
		var id int64
		var name string
		var javaPrice float64

		err := rows.Scan(&id, &name, &javaPrice)
		if err != nil {
			return fmt.Errorf("error scanning Java record: %v", err)
		}

		var goPrice int
		err = goDB.QueryRow("SELECT part_price FROM parts WHERE part_id = ?", id).Scan(&goPrice)
		if err != nil {
			return fmt.Errorf("error querying Go record: %v", err)
		}

		fmt.Printf("ID: %d, Name: %s\n", id, name)
		fmt.Printf("  Java price: %.2f\n", javaPrice)
		fmt.Printf("  Go price: %.2f\n", float64(goPrice)/100)
	}

	return nil
}
