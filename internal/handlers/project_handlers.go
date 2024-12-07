package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/device"
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
	id, err := upload.Int(e, "projectId")
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
	name := e.FormValue("projectName")

	p := project.NewProject(name)
	id, err := database.AddProjectToDB(p)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "project added successfully",
		"id":      strconv.Itoa(id),
	})
}

func CopyProject(e echo.Context) error {
	prjId, err := upload.Int(e, "projectId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	originalProject, err := database.GetProjectByIdFromDB(prjId)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	newProject := project.NewProject(originalProject.Name + " (Copy)")
	newProject.ProjectDeviceList = originalProject.ProjectDeviceList
	newProject.ExtraPriceList = originalProject.ExtraPriceList

	// Insure the name is unique
	for {
		err = database.CheckDeviceByNameFromDB(newProject.Name)
		if err == sql.ErrNoRows {
			break
		}
		newProject.Name += " (Copy)"
	}

	newId, err := database.AddProjectToDB(newProject)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// Add project devices
	for _, pd := range newProject.ProjectDeviceList {
		err = database.AddProjectDeviceToDB(newId, pd.Count, pd.Device.Id)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
	}

	// Add extra prices
	for _, ep := range newProject.ExtraPriceList {
		err = database.AddExtraPriceToDB(newId, ep.Name, ep.Price)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "project copied successfully",
		"id":      strconv.Itoa(newId),
	})
}

func DeleteProject(e echo.Context) error {
	prjId, err := upload.Int(e, "projectId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	err = database.DeleteProjectFromDB(prjId)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "project deleted successfully")
}

func UpdateProject(e echo.Context) error {
	id, err := upload.Int(e, "projectId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid project id")
	}

	name := e.FormValue("projectName")

	// Parse `devices` JSON from the form data
	devicesJSON := e.FormValue("devices")
	var devices []device.DeviceJson
	if err := json.Unmarshal([]byte(devicesJSON), &devices); err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	projectDevices := make([]device.DeviceReq, len(devices))
	for di, d := range devices {

		// convert deviceId to int
		deviceId, err := strconv.Atoi(d.Id)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}

		countf64, err := strconv.ParseFloat(d.Count, 32)
		if err != nil {
			return e.String(http.StatusInternalServerError, err.Error())
		}

		projectDevices[di] = device.DeviceReq{
			Id:    deviceId,
			Count: countf64,
		}
	}

	err = database.UpdateProjectInDB(id, name, projectDevices)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "project not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "project updated successfully")
}
