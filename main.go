package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lpar/gzipped"
	"github.com/phayes/freeport"
	embedded "github.com/poriamsz55/BoosterPump-webapp/embeded"
	"github.com/poriamsz55/BoosterPump-webapp/internal/lorca"
	"github.com/poriamsz55/BoosterPump-webapp/internal/routes"
)

func portSelection() (string, bool) {
	if len(os.Args) > 1 {
		port, _ := strconv.Atoi(os.Args[1])
		if port > 0 {
			return os.Args[1], true
		}
	}
	port, err := freeport.GetFreePort()
	if err != nil {
		fmt.Println("Could not find idle port", err)
		return strconv.Itoa(port), false
	}
	return strconv.Itoa(port), false
}

func WebServer(e embed.FS, routes ...func(e *echo.Echo)) string {
	port, block := portSelection()
	a := echo.New()
	for _, r := range routes {
		r(a)
	}

	f, err := fs.Sub(e, "public")
	if err != nil {
		fmt.Println("Could not embed front", err)
		return port
	}
	a.GET("/public/*", echo.WrapHandler(http.StripPrefix("/public/", gzipped.FileServer(http.FS(f)))))
	a.HideBanner = true
	opt := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	if block {
		a.StartServer(opt)
	}
	go a.StartServer(opt)
	return port
}

func main() {
	port := WebServer(embedded.BoosterFiles, routes.MainRoutes)
	lorca.Setup(port)
}
