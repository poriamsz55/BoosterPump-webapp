package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Add this method to your ExcelWriter struct
func (ew *ExcelWriter) SaveToDownloads(f *excelize.File, fileName string) (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	// Get Downloads directory based on OS
	var downloadsDir string
	switch runtime.GOOS {
	case "windows":
		downloadsDir = filepath.Join(homeDir, "Downloads")
	case "linux":
		downloadsDir = filepath.Join(homeDir, "Downloads")
		// Check if Downloads directory exists, if not try lowercase
		if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
			downloadsDir = filepath.Join(homeDir, "downloads")
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Create Downloads directory if it doesn't exist
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create Downloads directory: %v", err)
	}

	// Create a unique filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	safeName := sanitizeFileName(fileName)
	excelName := fmt.Sprintf("%s_%s.xlsx", safeName, timestamp)
	filePath := filepath.Join(downloadsDir, excelName)

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, append a number to the filename
		counter := 1
		for {
			newFileName := fmt.Sprintf("%s_%s_%d.xlsx", safeName, timestamp, counter)
			filePath = filepath.Join(downloadsDir, newFileName)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				break
			}
			counter++
		}
	}

	// Save the file
	err = f.SaveAs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save Excel file: %v", err)
	}

	return filePath, nil
}

// Helper function to sanitize file name
func sanitizeFileName(fileName string) string {
	// Replace invalid characters with underscores
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := fileName

	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Trim spaces from beginning and end
	result = strings.TrimSpace(result)

	// If filename is empty after sanitization, use a default name
	if result == "" {
		result = "project"
	}

	return result
}

// Optional: Add a method to check if Downloads directory is accessible
func (ew *ExcelWriter) CheckDownloadsAccess() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Get Downloads directory based on OS
	var downloadsDir string
	switch runtime.GOOS {
	case "windows":
		downloadsDir = filepath.Join(homeDir, "Downloads")
	case "linux":
		downloadsDir = filepath.Join(homeDir, "Downloads")
		// Check if Downloads directory exists, if not try lowercase
		if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
			downloadsDir = filepath.Join(homeDir, "downloads")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Try to create a test file
	testFile := filepath.Join(downloadsDir, ".test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("downloads directory is not writable: %v", err)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}
