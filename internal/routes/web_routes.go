package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func webRoutes(e *echo.Echo) {
	// main
	e.GET("/", handlers.IndexView)
	e.GET("/projects", handlers.ProjectsView)
	e.GET("/add/project/db", handlers.AddProjectDBView)
	e.GET("/devices", handlers.DevicesView)
	e.GET("/add/device/db", handlers.AddDeviceDBView)
	e.GET("/parts", handlers.PartsView)
	e.GET("/extra-prices", handlers.ExtraPricesView)

	// details
	e.GET("/projects/details", handlers.ProjectDetailsView)
	e.GET("/devices/details", handlers.DeviceDetailsView)
	e.GET("/parts/details", handlers.PartDetailsView)

}
