package loader

import (
	"database/sql"
	"fmt"
)

// Function to convert all tables (combining with previous code)
func ConvertAllTables(javaDB, goDB *sql.DB) error {
	// Convert projects
	if err := convertProjects(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting projects: %v", err)
	}

	// Convert extra prices
	if err := convertExtraPrices(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting extra prices: %v", err)
	}

	// Convert parts
	if err := convertParts(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting parts: %v", err)
	}

	// Convert devices
	if err := convertDevices(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting devices: %v", err)
	}

	// Convert project devices
	if err := convertProjectDevices(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting project devices: %v", err)
	}

	// Convert device parts
	if err := convertDeviceParts(javaDB, goDB); err != nil {
		return fmt.Errorf("error converting device parts: %v", err)
	}

	return nil
}

// Function to verify all conversions
func VerifyAllConversions(javaDB, goDB *sql.DB) error {
	// Verify projects
	if err := verifyProjectsConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("projects verification failed: %v", err)
	}

	// Verify devices
	if err := verifyDevicesConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("devices verification failed: %v", err)
	}

	// Verify parts
	if err := verifyPartsConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("parts verification failed: %v", err)
	}

	// Verify project devices
	if err := verifyProjectDevicesConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("project devices verification failed: %v", err)
	}

	// Verify extra prices
	if err := verifyExtraPricesConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("extra prices verification failed: %v", err)
	}

	// Verify device parts
	if err := verifyDevicePartsConversion(javaDB, goDB); err != nil {
		return fmt.Errorf("device parts verification failed: %v", err)
	}

	return nil
}
