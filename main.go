package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lpar/gzipped"
	embedded "github.com/poriamsz55/BoosterPump-webapp/embeded"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/lorca"
	"github.com/poriamsz55/BoosterPump-webapp/internal/routes"
	"github.com/poriamsz55/BoosterPump-webapp/internal/temp"
)

func WebServer(e embed.FS, port string, routes ...func(e *echo.Echo)) error {
	a := echo.New()

	// Enable CORS middleware
	a.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:8080", "http://localhost:8080"}, // Allowed origins
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},       // Allowed HTTP methods
	}))

	for _, r := range routes {
		r(a)
	}

	f, err := fs.Sub(e, "public")
	if err != nil {
		return err
	}
	a.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", gzipped.FileServer(http.FS(f)))))
	a.HideBanner = true
	opt := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	go a.StartServer(opt)
	return nil
}

func main() {

	// delete every temp database created
	// tempFile, err := os.CreateTemp("", "booster_db_temp_*.db")
	temp.CleanupTempDatabases()

	// start the database
	database.InitializeDB()
	defer database.CloseDB()

	port := "8080"
	WebServer(embedded.BoosterFiles, port, routes.MainRoutes)
	lorca.Setup(port)
}
