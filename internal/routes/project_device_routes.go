package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func projectDeviceRoutes(e *echo.Group) {
	e.GET("/projectDevice/getAll", handlers.GetAllProjectDevices)
	e.GET("/projectDevice/getById", handlers.GetProjectDeviceById)
	e.GET("/projectDevice/add", handlers.AddProjectDevice)
	e.GET("/projectDevice/copy", handlers.CopyProjectDevice)
	e.GET("/projectDevice/delete", handlers.DeleteProjectDevice)
	e.GET("/projectDevice/update", handlers.UpdateProjectDevice)
}
