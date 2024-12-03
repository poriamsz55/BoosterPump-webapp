package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers"
)

func projectRoutes(e *echo.Group) {
	e.GET("/project/getAll", handlers.GetAllProjects)
	e.POST("/project/getById", handlers.GetProjectById)
	e.POST("/project/add", handlers.AddProject)
	e.POST("/project/copy", handlers.CopyProject)
	e.POST("/project/delete", handlers.DeleteProject)
	e.POST("/project/update", handlers.UpdateProject)
}
