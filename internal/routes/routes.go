package routes

import (
	"github.com/labstack/echo/v4"
)

func MainRoutes(e *echo.Echo) {
	// Routes
	// web
	webRoutes(e)

	// api
	apiRoutes(e)

}

func apiRoutes(e *echo.Echo) {
	apiGrp := e.Group("/api")
	partRoutes(apiGrp)
	projectRoutes(apiGrp)
	deviceRoutes(apiGrp)
	extraPriceRoutes(apiGrp)
	projectDeviceRoutes(apiGrp)
	devicePartRoutes(apiGrp)
}
