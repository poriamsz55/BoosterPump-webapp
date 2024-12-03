package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func devicePartRoutes(e *echo.Group) {
	e.GET("/devicePart/getAll", handlers.GetAllDeviceParts)
	e.POST("/devicePart/getById", handlers.GetDevicePartById)
	e.POST("/devicePart/add", handlers.AddDevicePart)
	e.POST("/devicePart/add/list", handlers.AddDevicePartList)
	e.POST("/devicePart/copy", handlers.CopyDevicePart)
	e.POST("/devicePart/delete", handlers.DeleteDevicePart)
	e.POST("/devicePart/update", handlers.UpdateDevicePart)
}
