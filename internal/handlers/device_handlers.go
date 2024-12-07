package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func GetAllDevices(e echo.Context) error {
	devices, err := database.GetAllDevicesFromDB()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, devices)
}

func GetDeviceById(e echo.Context) error {
	id, err := upload.Int(e, "deviceId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device id")
	}

	device, err := getDeviceById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, device)
}

func getDeviceById(dvcId int) (*device.Device, error) {
	return database.GetDeviceByIdFromDB(dvcId)
}

func AddDevice(e echo.Context) error {
	name := e.FormValue("deviceName")
	converterInt, err := upload.Int(e, "converterType")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	converter, err := device.ConverterFromValue(converterInt)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	filter := e.FormValue("filter") == "true"

	d := device.NewDevice(name, converter, filter)
	id, err := database.AddDeviceToDB(d)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "device added successfully",
		"id":      strconv.Itoa(id),
	})
}

func CopyDevice(e echo.Context) error {
	id, err := upload.Int(e, "deviceId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device id")
	}

	originalDevice, err := getDeviceById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// Create new device with copied properties
	newDevice := device.NewDevice(
		fmt.Sprintf("%s (Copy)", originalDevice.Name),
		originalDevice.Converter,
		originalDevice.Filter,
	)
	newDevice.DevicePartList = originalDevice.DevicePartList

	// Insure the name is unique
	for {
		err = database.CheckDeviceByNameFromDB(newDevice.Name)
		if err == sql.ErrNoRows {
			break
		}
		newDevice.Name += " (Copy)"
	}

	copyId, err := database.AddDeviceToDB(newDevice)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// Copy device parts
	for _, dp := range newDevice.DevicePartList {
		err = database.AddDevicePartToDB(copyId, dp.Count, dp.Part.Id)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "device copied successfully",
		"id":      strconv.Itoa(copyId),
	})
}

func DeleteDevice(e echo.Context) error {
	id, err := upload.Int(e, "deviceId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device id")
	}

	err = database.DeleteDeviceFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device deleted successfully")
}

func UpdateDevice(e echo.Context) error {
	id, err := upload.Int(e, "deviceId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device id")
	}

	name := e.FormValue("deviceName")
	converterInt, err := upload.Int(e, "converterType")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	converter, err := device.ConverterFromValue(converterInt)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	filter := e.FormValue("filter") == "true"

	updatedDevice := device.NewDevice(name, converter, filter)
	updatedDevice.Id = id

	// Parse `parts` JSON from the form data
	partsJSON := e.FormValue("parts")
	var parts []part.PartJson
	if err := json.Unmarshal([]byte(partsJSON), &parts); err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	deviceParts := make([]part.PartReq, len(parts))
	for pi, p := range parts {

		// convert partId to int
		partId, err := strconv.Atoi(p.Id)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}

		countf64, err := strconv.ParseFloat(p.Count, 32)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}

		deviceParts[pi] = part.PartReq{
			Id:    partId,
			Count: countf64,
		}
	}

	err = database.UpdateDeviceInDB(updatedDevice, deviceParts)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device updated successfully")
}
