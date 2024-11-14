package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	embedded "github.com/poriamsz55/BoosterPump-webapp/embeded"
)

func MainRoutes(e *echo.Echo) {
	// Routes
	e.GET("/", handleIndex)
	e.GET("/projects", handleProjects)
	e.GET("/distork", handleDistork)
}

// Handlers
func handleIndex(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/index.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func handleProjects(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/projects.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func handleDistork(c echo.Context) error {
	html, err := embedded.BoosterFiles.ReadFile("views/distork.html")
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
