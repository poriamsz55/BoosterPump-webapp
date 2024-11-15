package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	embedded "github.com/poriamsz55/BoosterPump-webapp/embeded"
)

// Handlers
func HandleIndex(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/index.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func HandleProjects(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/projects.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func HandleDevices(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/devices.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func HandleParts(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/parts.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Detail

func HandleProjectDetails(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/project-details.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func HandleDeviceDetails(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/device-details.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func HandlePartDetails(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/part-details.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
