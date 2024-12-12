package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
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
	databaseRoutes(apiGrp)
	partRoutes(apiGrp)
	projectRoutes(apiGrp)
	deviceRoutes(apiGrp)
	extraPriceRoutes(apiGrp)
	projectDeviceRoutes(apiGrp)
	devicePartRoutes(apiGrp)
}

func databaseRoutes(apiGrp *echo.Group) {
	apiGrp.POST("/database/download", handlers.DownloadDatabase)
	apiGrp.POST("/database/upload", handlers.UploadDatabase)
}
