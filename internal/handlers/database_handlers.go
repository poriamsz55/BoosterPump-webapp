package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/temp"
)

func DownloadDatabase(e echo.Context) error {
	saveName := e.FormValue("name")
	if saveName == "" {
		saveName = "booster_pump.db"
	}

	// Ensure the filename is safe
	saveName = strings.ReplaceAll(saveName, "/", "")
	saveName = strings.ReplaceAll(saveName, "\\", "")

	// Ensure it has .db extension
	if !strings.HasSuffix(saveName, ".db") {
		saveName += ".db"
	}

	// Check if database file exists
	if _, err := os.Stat("booster_pump.db"); os.IsNotExist(err) {
		return e.String(http.StatusNotFound, "Database file not found")
	}

	// Open the database file
	dbFile, err := os.Open("booster_pump.db")
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error opening database file: "+err.Error())
	}
	defer dbFile.Close()

	// Get file information
	fileInfo, err := dbFile.Stat()
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error getting file info: "+err.Error())
	}

	// Set response headers
	e.Response().Header().Set("Content-Type", "application/octet-stream")
	e.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", saveName))
	e.Response().Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Copy the file to the response writer
	_, err = io.Copy(e.Response().Writer, dbFile)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error sending file: "+err.Error())
	}

	return nil
}

func UploadDatabase(e echo.Context) error {
	temp.CleanupTempDatabases()

	// Get the uploaded file
	file, err := e.FormFile("database")
	if err != nil {
		return e.String(http.StatusBadRequest, "Error getting uploaded file: "+err.Error())
	}

	// Get replace value
	replace := e.FormValue("replace") == "true"

	// Create a temporary file to store the upload
	tempFile, err := os.CreateTemp("", "booster_db_temp_*.db")
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error creating temp file: "+err.Error())
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error opening uploaded file: "+err.Error())
	}
	defer src.Close()

	// Copy uploaded file to temp file
	_, err = io.Copy(tempFile, src)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error copying to temp file: "+err.Error())
	}

	if replace {
		database.CloseDB()

		// Remove existing database
		if err = os.Remove("booster_pump.db"); err != nil && !os.IsNotExist(err) {
			return e.String(http.StatusInternalServerError, "Error removing existing database: "+err.Error())
		}

		// Copy temp file to database file
		dst, err := os.Create("booster_pump.db")
		if err != nil {
			return e.String(http.StatusInternalServerError, "Error creating new database: "+err.Error())
		}
		defer dst.Close()

		if _, err = tempFile.Seek(0, 0); err != nil {
			return e.String(http.StatusInternalServerError, "Error seeking temp file: "+err.Error())
		}

		if _, err = io.Copy(dst, tempFile); err != nil {
			return e.String(http.StatusInternalServerError, "Error copying new database: "+err.Error())
		}

		database.RunDB()
	} else {
		// Merge databases
		err := database.Merge(tempFile.Name())
		if err != nil {
			return e.String(http.StatusInternalServerError, "Error merging database: "+err.Error())
		}
	}

	return e.String(http.StatusOK, "Database uploaded successfully")
}
