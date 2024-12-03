package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func partRoutes(e *echo.Group) {
	e.GET("/part/getAll", handlers.GetAllParts)
	e.POST("/part/getById", handlers.GetPartById)
	e.POST("/part/add", handlers.AddPart)
	e.POST("/part/copy", handlers.CopyPart)
	e.POST("/part/delete", handlers.DeletePart)
	e.POST("/part/update", handlers.UpdatePart)
}
