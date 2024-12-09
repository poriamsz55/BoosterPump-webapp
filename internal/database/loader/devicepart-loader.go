package loader

import (
	"database/sql"
	"fmt"
)

// Struct for Java DevicePart (source)
type JavaDevicePart struct {
	ID       int64
	Count    float64 // FLOAT in Java
	PartID   int64
	DeviceID int64
}

// Struct for Go DevicePart (destination)
type GoDevicePart struct {
	ID       int64
	Count    float64 // Keep as FLOAT in Go as per your schema
	PartID   int64
	DeviceID int64
}

func convertDeviceParts(javaDB, goDB *sql.DB) error {
	// Query all device parts from Java DB
	rows, err := javaDB.Query(`
        SELECT device_part_id, device_part_count, part_id, device_id 
        FROM device_parts
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
        INSERT INTO device_parts (
            device_part_id,
            device_part_count,
            part_id,
            device_id
        ) VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Process each device part
	var count int
	for rows.Next() {
		var jDevicePart JavaDevicePart
		err := rows.Scan(
			&jDevicePart.ID,
			&jDevicePart.Count,
			&jDevicePart.PartID,
			&jDevicePart.DeviceID,
		)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Convert Java device part to Go device part
		gDevicePart := convertJavaToGoDevicePart(jDevicePart)

		// Insert into Go database
		_, err = stmt.Exec(
			gDevicePart.ID,
			gDevicePart.Count,
			gDevicePart.PartID,
			gDevicePart.DeviceID,
		)
		if err != nil {
			return fmt.Errorf("error inserting device part %d: %v", gDevicePart.ID, err)
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

	fmt.Printf("Converted %d device parts\n", count)
	return nil
}

func convertJavaToGoDevicePart(jDevicePart JavaDevicePart) GoDevicePart {
	return GoDevicePart{
		ID:       jDevicePart.ID,
		Count:    jDevicePart.Count, // Keep as float64
		PartID:   jDevicePart.PartID,
		DeviceID: jDevicePart.DeviceID,
	}
}

func verifyDevicePartsConversion(javaDB, goDB *sql.DB) error {
	var javaCount, goCount int

	// Count records in Java DB
	err := javaDB.QueryRow("SELECT COUNT(*) FROM device_parts").Scan(&javaCount)
	if err != nil {
		return fmt.Errorf("error counting Java records: %v", err)
	}

	// Count records in Go DB
	err = goDB.QueryRow("SELECT COUNT(*) FROM device_parts").Scan(&goCount)
	if err != nil {
		return fmt.Errorf("error counting Go records: %v", err)
	}

	fmt.Printf("Device parts count - Java DB: %d, Go DB: %d\n", javaCount, goCount)

	if javaCount != goCount {
		return fmt.Errorf("record count mismatch: Java=%d, Go=%d", javaCount, goCount)
	}

	// Verify some sample records with related data
	rows, err := javaDB.Query(`
        SELECT 
            dp.device_part_id,
            dp.device_part_count,
            p.part_name,
            d.device_name
        FROM device_parts dp
        JOIN parts p ON dp.part_id = p.part_id
        JOIN devices d ON dp.device_id = d.device_id
        LIMIT 5
    `)
	if err != nil {
		return fmt.Errorf("error querying sample records: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nSample device parts comparison:")
	for rows.Next() {
		var id int64
		var count float64
		var partName, deviceName string

		err := rows.Scan(&id, &count, &partName, &deviceName)
		if err != nil {
			return fmt.Errorf("error scanning Java record: %v", err)
		}

		var goCount float64
		err = goDB.QueryRow(`
            SELECT dp.device_part_count 
            FROM device_parts dp 
            WHERE dp.device_part_id = ?`, id).Scan(&goCount)
		if err != nil {
			return fmt.Errorf("error querying Go record: %v", err)
		}

		fmt.Printf("ID: %d, Part: %s, Device: %s\n", id, partName, deviceName)
		fmt.Printf("  Java count: %.2f\n", count)
		fmt.Printf("  Go count: %.2f\n", goCount)
	}

	return nil
}
