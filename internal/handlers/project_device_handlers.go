package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
)

func GetAllProjectDevices(e echo.Context) error {
	panic("to be implement")
}

func GetProjectDeviceById(e echo.Context) error {
	panic("to be implement")
}

func AddProjectDevice(e echo.Context) error {
	prjId, err := upload.Int(e, "project_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	count, err := upload.Float64(e, "count")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	dvcId, err := upload.Int(e, "device_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	err = database.AddProjectDeviceToDB(prjId, count, dvcId)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "projectDevice added to database successfully")
}

type DeviceJson struct {
	Id    string `json:"id"`
	Count string `json:"count"`
}

func AddProjectDeviceList(e echo.Context) error {
	prjId, err := upload.Int(e, "projectId")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// Parse `devices` JSON from the form data
	devicesJSON := e.FormValue("devices")
	var devices []DeviceJson
	if err := json.Unmarshal([]byte(devicesJSON), &devices); err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	for _, d := range devices {

		// convert deviceId to int
		deviceId, err := strconv.Atoi(d.Id)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}

		countf64, err := strconv.ParseFloat(d.Count, 32)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
		err = database.AddProjectDeviceToDB(prjId, float64(countf64), deviceId)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
	}

	return e.String(http.StatusOK, "projectDevice added to database successfully")
}

func CopyProjectDevice(e echo.Context) error {
	panic("to be implement")
}

func DeleteProjectDevice(e echo.Context) error {
	panic("to be implement")
}

func UpdateProjectDevice(e echo.Context) error {
	panic("to be implement")
}
