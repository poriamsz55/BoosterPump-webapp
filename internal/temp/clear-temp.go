package temp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CleanupTempDatabases() error {
	// Get the system's temporary directory
	tempDir := os.TempDir()

	// Read all files in temp directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("error reading temp directory: %v", err)
	}

	// Look for our temp database files and remove them
	for _, entry := range entries {
		if entry.Type().IsRegular() && // is a regular file
			strings.HasPrefix(entry.Name(), "booster_db_temp_") && // has our prefix
			strings.HasSuffix(entry.Name(), ".db") { // has .db extension

			fullPath := filepath.Join(tempDir, entry.Name())
			err := os.Remove(fullPath)
			if err != nil {
				// Log the error but continue with other files
				fmt.Printf("Warning: Failed to delete temp file %s: %v\n", fullPath, err)
			}
		}
	}

	return nil
}
