package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func extraPriceRoutes(e *echo.Group) {
	e.GET("/extraPrice/getAll", handlers.GetAllExtraPrices)
	e.GET("/extraPrice/getById", handlers.GetExtraPriceById)
	e.GET("/extraPrice/add", handlers.AddExtraPrice)
	e.GET("/extraPrice/copy", handlers.CopyExtraPrice)
	e.GET("/extraPrice/delete", handlers.DeleteExtraPrice)
	e.GET("/extraPrice/update", handlers.UpdateExtraPrice)
}
