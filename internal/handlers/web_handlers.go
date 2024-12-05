package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	embedded "github.com/poriamsz55/BoosterPump-webapp/embeded"
)

// Handlers
func IndexView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/index.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func ProjectsView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/projects.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func AddProjectDBView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/project-db.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func ExtraPricesView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/extra-prices.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func DevicesView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/devices.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func AddDeviceDBView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/device-db.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func PartsView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/parts.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

// Detail

func ProjectDetailsView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/project-details.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func DeviceDetailsView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/device-details.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func PartDetailsView(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/part-details.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
