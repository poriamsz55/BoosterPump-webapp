package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func webRoutes(e *echo.Echo) {
	// main
	e.GET("/", handlers.HandleIndex)
	e.GET("/projects", handlers.HandleProjects)
	e.GET("/devices", handlers.HandleDevices)
	e.GET("/parts", handlers.HandleParts)

	// details
	e.GET("/projects/details", handlers.HandleProjectDetails)
	e.GET("/devices/details", handlers.HandleDeviceDetails)
	e.GET("/parts/details", handlers.HandlePartDetails)

}
