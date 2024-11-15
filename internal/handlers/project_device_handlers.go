package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
)

func GetAllProjectDevices(e echo.Context) error {
	panic("dasdas")
}

func GetProjectDeviceById(e echo.Context) error {
	panic("dasdas")
}

func AddProjectDevice(e echo.Context) error {
	prjId, err := upload.Int(e, "project_id")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	count, err := upload.Float32(e, "count")
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

func CopyProjectDevice(e echo.Context) error {
	panic("dasdas")
}

func DeleteProjectDevice(e echo.Context) error {
	panic("dasdas")
}

func UpdateProjectDevice(e echo.Context) error {
	panic("dasdas")
}
