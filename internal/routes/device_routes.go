package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func deviceRoutes(e *echo.Group) {
	e.GET("/device/getAll", handlers.GetAllDevices)
	e.POST("/device/getById", handlers.GetDeviceById)
	e.POST("/device/add", handlers.AddDevice)
	e.POST("/device/copy", handlers.CopyDevice)
	e.POST("/device/delete", handlers.DeleteDevice)
	e.POST("/device/update", handlers.UpdateDevice)
}
