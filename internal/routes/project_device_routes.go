package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func projectDeviceRoutes(e *echo.Group) {
	e.GET("/projectDevice/getAll", handlers.GetAllProjectDevices)
	e.POST("/projectDevice/getById", handlers.GetProjectDeviceById)
	e.POST("/projectDevice/add", handlers.AddProjectDevice)
	e.POST("/projectDevice/copy", handlers.CopyProjectDevice)
	e.POST("/projectDevice/delete", handlers.DeleteProjectDevice)
	e.POST("/projectDevice/update", handlers.UpdateProjectDevice)
}
