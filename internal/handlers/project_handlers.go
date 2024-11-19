package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/project"
)

func GetAllProjects(e echo.Context) error {
	projects, err := database.GetAllProjectsFromDB()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, projects)
}

func GetProjectById(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	project, err := database.GetProjectByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, project)
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
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	originalProject, err := database.GetProjectByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	newProject := project.NewProject(originalProject.Name + " (Copy)")
	newProject.ProjectDeviceList = originalProject.ProjectDeviceList
	newProject.ExtraPriceList = originalProject.ExtraPriceList

	err = database.AddProjectToDB(newProject)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "project copied successfully")
}

func DeleteProject(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	err = database.DeleteProjectFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "project deleted successfully")
}

func UpdateProject(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	name := e.FormValue("name")

	err = database.UpdateProjectInDB(id, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "project updated successfully")
}
