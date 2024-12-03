package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
)

func GetAllDevices(e echo.Context) error {
	devices, err := database.GetAllDevicesFromDB()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, devices)
}

func GetDeviceById(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
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
	id, err := strconv.Atoi(e.Param("id"))
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

	copyId, err := database.AddDeviceToDB(newDevice)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "device copied successfully",
		"id":      strconv.Itoa(copyId),
	})
}

func DeleteDevice(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
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
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid device id")
	}

	name := e.FormValue("name")
	converterInt, err := upload.Int(e, "converter")
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

	err = database.UpdateDeviceInDB(updatedDevice)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "device not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device updated successfully")
}
