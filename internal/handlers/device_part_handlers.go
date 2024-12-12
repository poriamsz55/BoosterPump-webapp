package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func GetAllDeviceParts(e echo.Context) error {
	deviceId, err := strconv.Atoi(e.QueryParam("device_id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device id")
	}

	deviceParts, err := database.GetDevicePartsByDeviceId(nil, deviceId)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, deviceParts)
}

func GetDevicePartById(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device part id")
	}

	devicePart, err := database.GetDevicePartByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, devicePart)
}

func AddDevicePart(e echo.Context) error {
	deviceId, err := upload.Int(e, "device_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	count, err := upload.Float64(e, "count")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	partId, err := upload.Int(e, "part_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	err = database.AddDevicePartToDB(deviceId, count, partId)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "devicePart added to database successfully")
}

func AddDevicePartList(e echo.Context) error {

	// Parse `deviceId` from the form data
	deviceIDStr := e.FormValue("deviceId")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// Parse `parts` JSON from the form data
	partsJSON := e.FormValue("parts")
	var parts []part.PartJson
	if err := json.Unmarshal([]byte(partsJSON), &parts); err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	for _, p := range parts {

		// convert partId to int
		partId, err := strconv.Atoi(p.Id)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}

		countf64, err := strconv.ParseFloat(p.Count, 32)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
		err = database.AddDevicePartToDB(deviceID, float64(countf64), partId)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
	}

	return e.String(http.StatusOK, "deviceParts added to database successfully")
}

func CopyDevicePart(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device part id")
	}

	originalDevicePart, err := database.GetDevicePartByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	err = database.AddDevicePartToDB(
		originalDevicePart.DeviceId,
		originalDevicePart.Count,
		originalDevicePart.Part.Id,
	)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device part copied successfully")
}

func DeleteDevicePart(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device part id")
	}

	err = database.DeleteDevicePartFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device part deleted successfully")
}

func UpdateDevicePart(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device part id")
	}

	count, err := upload.Float64(e, "count")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	deviceId, err := upload.Int(e, "device_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	projectId, err := upload.Int(e, "project_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	err = database.UpdateDevicePartInDB(id, projectId, count, deviceId)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device part updated successfully")
}
