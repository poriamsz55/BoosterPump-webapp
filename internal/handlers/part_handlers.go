package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func GetAllParts(e echo.Context) error {
	parts, err := database.GetAllPartsFromDB()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, parts)
}

func GetPartById(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid part id")
	}

	part, err := database.GetPartByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, part)
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
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid part id")
	}

	originalPart, err := database.GetPartByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	newPart := part.NewPart(
		originalPart.Name+" (Copy)",
		originalPart.Size,
		originalPart.Material,
		originalPart.Brand,
		originalPart.Price,
	)

	err = database.AddPartToDB(newPart)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "part copied successfully")
}

func DeletePart(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid part id")
	}

	err = database.DeletePartFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "part deleted successfully")
}

func UpdatePart(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid part id")
	}

	name := e.FormValue("name")
	size := e.FormValue("size")
	material := e.FormValue("material")
	brand := e.FormValue("brand")
	price, err := upload.Uint64(e, "price")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	updatedPart := part.NewPart(name, size, material, brand, price)
	updatedPart.Id = id

	err = database.UpdatePartInDB(updatedPart)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "part not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "part updated successfully")
}
