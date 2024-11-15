package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func projectRoutes(e *echo.Group) {
	e.GET("/project/getAll", handlers.GetAllProjects)
	e.GET("/project/getById", handlers.GetProjectById)
	e.GET("/project/add", handlers.AddProject)
	e.GET("/project/copy", handlers.CopyProject)
	e.GET("/project/delete", handlers.DeleteProject)
	e.GET("/project/update", handlers.UpdateProject)
}
