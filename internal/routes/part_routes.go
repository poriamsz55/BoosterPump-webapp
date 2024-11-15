package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func partRoutes(e *echo.Group) {
	e.GET("/part/getAll", handlers.GetAllParts)
	e.GET("/part/getById", handlers.GetPartById)
	e.GET("/part/add", handlers.AddPart)
	e.GET("/part/copy", handlers.CopyPart)
	e.GET("/part/delete", handlers.DeletePart)
	e.GET("/part/update", handlers.UpdatePart)
}
