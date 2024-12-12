package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
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

	// Get the file path form the request
	path := e.FormValue("path")

	// get replace value
	replace := e.FormValue("replace") == "true"

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error opening file: "+err.Error())
	}
	defer file.Close()

	// if replace is true, replace new database file with the database file
	if replace {
		err := os.Remove("booster_pump.db")
		if err != nil {
			return e.String(http.StatusInternalServerError, "Error removing database file: "+err.Error())
		}

		// Copy the file to the database file
		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			return e.String(http.StatusInternalServerError, "Error copying file: "+err.Error())
		}
	} else {
		// if !replace, add new database rows to the database file
		err := database.Merge(path)
		if err != nil {
			return e.String(http.StatusInternalServerError, "Error merging database: "+err.Error())
		}
	}

	return e.String(http.StatusOK, "Database file uploaded successfully")
}
