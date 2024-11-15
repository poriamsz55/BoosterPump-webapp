package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
)

func GetAllDevices(e echo.Context) error {
	panic("dasdas")
}

func GetDeviceById(e echo.Context) error {
	panic("dasdas")
}

func getDeviceById(dvcId int) (*device.Device, error) {
	panic("dasdas")
}

func AddDevice(e echo.Context) error {
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

	d := device.NewDevice(name, converter, filter)
	err = database.AddDeviceToDB(d)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "device added to database successfully")
}

func CopyDevice(e echo.Context) error {
	panic("dasdas")
}

func DeleteDevice(e echo.Context) error {
	panic("dasdas")
}

func UpdateDevice(e echo.Context) error {
	panic("dasdas")
}
