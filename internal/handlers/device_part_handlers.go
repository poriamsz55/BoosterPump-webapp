package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
)

func GetAllDeviceParts(e echo.Context) error {
	panic("dasdas")
}

func GetDevicePartById(e echo.Context) error {
	panic("dasdas")
}

func AddDevicePart(e echo.Context) error {
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

	err = database.AddDevicePartToDB(prjId, count, dvcId)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "devicePart added to database successfully")
}

func CopyDevicePart(e echo.Context) error {
	panic("dasdas")
}

func DeleteDevicePart(e echo.Context) error {
	panic("dasdas")
}

func UpdateDevicePart(e echo.Context) error {
	panic("dasdas")
}
