package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func GetAllParts(e echo.Context) error {
	panic("dasdas")
}

func GetPartById(e echo.Context) error {
	panic("dasdas")
}

func AddPart(e echo.Context) error {
	name := e.FormValue("name")
	size := e.FormValue("size")
	material := e.FormValue("material")
	brand := e.FormValue("brand")
	price, err := upload.Uint64(e, "price")

	p := part.NewPart(name, size, material, brand, price)
	err = database.AddPartToDB(p)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "part added to database successfully")
}

func CopyPart(e echo.Context) error {
	panic("dasdas")
}

func DeletePart(e echo.Context) error {
	panic("dasdas")
}

func UpdatePart(e echo.Context) error {
	panic("dasdas")
}
