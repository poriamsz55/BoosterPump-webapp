package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func devicePartRoutes(e *echo.Group) {
	e.GET("/devicePart/getAll", handlers.GetAllDeviceParts)
	e.GET("/devicePart/getById", handlers.GetDevicePartById)
	e.GET("/devicePart/add", handlers.AddDevicePart)
	e.GET("/devicePart/copy", handlers.CopyDevicePart)
	e.GET("/devicePart/delete", handlers.DeleteDevicePart)
	e.GET("/devicePart/update", handlers.UpdateDevicePart)
}
