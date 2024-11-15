package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func deviceRoutes(e *echo.Group) {
	e.GET("/device/getAll", handlers.GetAllDevices)
	e.GET("/device/getById", handlers.GetDeviceById)
	e.GET("/device/add", handlers.AddDevice)
	e.GET("/device/copy", handlers.CopyDevice)
	e.GET("/device/delete", handlers.DeleteDevice)
	e.GET("/device/update", handlers.UpdateDevice)
}
