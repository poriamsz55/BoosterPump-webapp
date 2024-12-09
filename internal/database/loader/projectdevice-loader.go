package loader

import (
	"database/sql"
	"fmt"
)

// Struct for Java ProjectDevice (source)
type JavaProjectDevice struct {
	ID        int64
	Count     float64 // FLOAT in Java
	DeviceID  int64
	ProjectID int64
}

// Struct for Go ProjectDevice (destination)
type GoProjectDevice struct {
	ID        int64
	Count     float64 // Keep as FLOAT in Go as per your schema
	DeviceID  int64
	ProjectID int64
}

func convertProjectDevices(javaDB, goDB *sql.DB) error {
	// Query all project devices from Java DB
	rows, err := javaDB.Query(`
        SELECT project_device_id, project_device_count, device_id, project_id 
        FROM project_devices
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
        INSERT INTO project_devices (
            project_device_id,
            project_device_count,
            device_id,
            project_id
        ) VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Process each project device
	var count int
	for rows.Next() {
		var jProjectDevice JavaProjectDevice
		err := rows.Scan(
			&jProjectDevice.ID,
			&jProjectDevice.Count,
			&jProjectDevice.DeviceID,
			&jProjectDevice.ProjectID,
		)
		if err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		// Convert Java project device to Go project device
		gProjectDevice := convertJavaToGoProjectDevice(jProjectDevice)

		// Insert into Go database
		_, err = stmt.Exec(
			gProjectDevice.ID,
			gProjectDevice.Count,
			gProjectDevice.DeviceID,
			gProjectDevice.ProjectID,
		)
		if err != nil {
			return fmt.Errorf("error inserting project device %d: %v", gProjectDevice.ID, err)
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

	fmt.Printf("Converted %d project devices\n", count)
	return nil
}

func convertJavaToGoProjectDevice(jProjectDevice JavaProjectDevice) GoProjectDevice {
	return GoProjectDevice{
		ID:        jProjectDevice.ID,
		Count:     jProjectDevice.Count, // Keep as float64
		DeviceID:  jProjectDevice.DeviceID,
		ProjectID: jProjectDevice.ProjectID,
	}
}

func verifyProjectDevicesConversion(javaDB, goDB *sql.DB) error {
	var javaCount, goCount int

	// Count records in Java DB
	err := javaDB.QueryRow("SELECT COUNT(*) FROM project_devices").Scan(&javaCount)
	if err != nil {
		return fmt.Errorf("error counting Java records: %v", err)
	}

	// Count records in Go DB
	err = goDB.QueryRow("SELECT COUNT(*) FROM project_devices").Scan(&goCount)
	if err != nil {
		return fmt.Errorf("error counting Go records: %v", err)
	}

	fmt.Printf("Project devices count - Java DB: %d, Go DB: %d\n", javaCount, goCount)

	if javaCount != goCount {
		return fmt.Errorf("record count mismatch: Java=%d, Go=%d", javaCount, goCount)
	}

	// Verify some sample records with related data
	rows, err := javaDB.Query(`
        SELECT 
            pd.project_device_id,
            pd.project_device_count,
            d.device_name,
            p.project_name
        FROM project_devices pd
        JOIN devices d ON pd.device_id = d.device_id
        JOIN projects p ON pd.project_id = p.project_id
        LIMIT 5
    `)
	if err != nil {
		return fmt.Errorf("error querying sample records: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nSample project devices comparison:")
	for rows.Next() {
		var id int64
		var count float64
		var deviceName, projectName string

		err := rows.Scan(&id, &count, &deviceName, &projectName)
		if err != nil {
			return fmt.Errorf("error scanning Java record: %v", err)
		}

		var goCount float64
		err = goDB.QueryRow(`
            SELECT pd.project_device_count 
            FROM project_devices pd 
            WHERE pd.project_device_id = ?`, id).Scan(&goCount)
		if err != nil {
			return fmt.Errorf("error querying Go record: %v", err)
		}

		fmt.Printf("ID: %d, Device: %s, Project: %s\n", id, deviceName, projectName)
		fmt.Printf("  Java count: %.2f\n", count)
		fmt.Printf("  Go count: %.2f\n", goCount)
	}

	return nil
}

// Update the convertAllTables function to include project devices
func convertAllTables(javaDB, goDB *sql.DB) error {
	// Convert in order of dependencies

	// 1. First convert projects (no dependencies)
	if err := convertProjects(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting projects: %v", err)
	}

	// 2. Convert devices (no dependencies)
	if err := convertDevices(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting devices: %v", err)
	}

	// 3. Convert project devices (depends on both projects and devices)
	if err := convertProjectDevices(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting project devices: %v", err)
	}

	// 4. Convert extra prices (depends on projects)
	if err := convertExtraPrices(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting extra prices: %v", err)
	}

	// 5. Convert parts (no dependencies)
	if err := convertParts(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting parts: %v", err)
	}

	return nil
}

// Update the verifyAllConversions function
func verifyAllConversions(javaDB, goDB *sql.DB) error {
	if err := verifyProjectsConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("projects verification failed: %v", err)
	}

	if err := verifyDevicesConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("devices verification failed: %v", err)
	}

	if err := verifyProjectDevicesConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("project devices verification failed: %v", err)
	}

	if err := verifyExtraPricesConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("extra prices verification failed: %v", err)
	}

	if err := verifyPartsConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("parts verification failed: %v", err)
	}

	return nil
}
