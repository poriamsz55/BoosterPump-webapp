package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func extraPriceRoutes(e *echo.Group) {
	e.GET("/extraPrice/getAll", handlers.GetAllExtraPrices)
	e.POST("/extraPrice/getById", handlers.GetExtraPriceById)
	e.POST("/extraPrice/add", handlers.AddExtraPrice)
	e.POST("/extraPrice/copy", handlers.CopyExtraPrice)
	e.POST("/extraPrice/delete", handlers.DeleteExtraPrice)
	e.POST("/extraPrice/update", handlers.UpdateExtraPrice)
}
