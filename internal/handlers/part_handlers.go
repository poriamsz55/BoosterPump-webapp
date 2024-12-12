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
	parts, err := database.GetAllParts(nil)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, parts)
}

func GetPartById(e echo.Context) error {
	id, err := strconv.Atoi(e.FormValue("partId"))
	if err != nil {
		e.JSON(http.StatusBadRequest, "invalid part id")
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
	name := e.FormValue("partName")
	size := e.FormValue("partSize")
	material := e.FormValue("partMaterial")
	brand := e.FormValue("partBrand")
	price, err := upload.Uint64(e, "partPrice")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	p := part.NewPart(name, size, material, brand, price)
	id, err := database.AddPartToDB(p)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "part added successfully",
		"id":      strconv.Itoa(id),
	})
}

func CopyPart(e echo.Context) error {
	id, err := upload.Int(e, "partId")
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

	// Insure the name is unique
	for {
		err = database.CheckPartByNameFromDB(newPart.Name)
		if err == sql.ErrNoRows {
			break
		}
		newPart.Name += " (Copy)"
	}

	newpId, err := database.AddPartToDB(newPart)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]string{
		"message": "part copied successfully",
		"id":      strconv.Itoa(newpId),
	})
}

func DeletePart(e echo.Context) error {
	id, err := upload.Int(e, "partId")
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
	id, err := strconv.Atoi(e.FormValue("partId"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid part id")
	}

	name := e.FormValue("partName")
	size := e.FormValue("partSize")
	material := e.FormValue("partMaterial")
	brand := e.FormValue("partBrand")
	price, err := upload.Uint64(e, "partPrice")
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
