package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/database"
)

func main() {
	database.InitializeDB()

	e := echo.New()

	// Serve static files
	e.Static("/", "web")

	// Splash screen route
	e.GET("/api", splashHandler)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

func splashHandler(c echo.Context) error {

	// Render the splash screen template
	return c.Render(http.StatusOK, "/web/index.html", nil)
}
