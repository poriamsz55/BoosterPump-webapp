package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
)

func GetAllProjects(e echo.Context) error {
	panic("dasdas")
}

func GetProjectById(e echo.Context) error {
	panic("dasdas")
}

func AddProject(e echo.Context) error {
	name := e.FormValue("name")

	p := project.NewProject(name)
	err := database.AddProjectToDB(p)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "project added to database successfully")
}

func CopyProject(e echo.Context) error {
	panic("dasdas")
}

func DeleteProject(e echo.Context) error {
	panic("dasdas")
}

func UpdateProject(e echo.Context) error {
	panic("dasdas")
}
