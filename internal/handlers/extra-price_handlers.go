package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database"
	"github.com/poriamsz55/BoosterPump-webapp/internal/handlers/upload"
	extraprice "github.com/poriamsz55/BoosterPump-webapp/internal/models/extra_price"
)

func GetAllExtraPricesByProjectId(e echo.Context) error {
	projectId, err := upload.Int(e, "projectId")
	extraPrices, err := database.GetExtraPricesByProjectIdFromDB(projectId)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, extraPrices)
}

func GetExtraPriceById(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid extraPrice id")
	}

	extraPrice, err := database.GetExtraPriceByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "extraPrice not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, extraPrice)
}

func AddExtraPrice(e echo.Context) error {
	name := e.FormValue("extraPriceName")
	projectId, err := upload.Int(e, "projectId")
	price, err := upload.Uint64(e, "extraPriceValue")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	p := extraprice.NewExtraPrice(projectId, name, price)
	err = database.AddExtraPriceToDB(p)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "extraPrice added to database successfully")
}

func CopyExtraPrice(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid extraPrice id")
	}

	originalExtraPrice, err := database.GetExtraPriceByIdFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "extraPrice not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	newExtraPrice := extraprice.NewExtraPrice(
		originalExtraPrice.ProjectId,
		originalExtraPrice.Name+" (Copy)",
		originalExtraPrice.Price,
	)

	// check if extraPrice already exists

	err = database.AddExtraPriceToDB(newExtraPrice)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "extraPrice copied successfully")
}

func DeleteExtraPrice(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid extraPrice id")
	}

	err = database.DeleteExtraPriceFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "extraPrice not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "extraPrice deleted successfully")
}

func UpdateExtraPrice(e echo.Context) error {
	id, err := upload.Int(e, "extraPriceId")
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid extraPrice id")
	}

	name := e.FormValue("extraPriceName")
	price, err := upload.Uint64(e, "extraPriceValue")
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	updatedExtraPrice := extraprice.NewExtraPrice(-1, name, price)
	updatedExtraPrice.Id = id

	err = database.UpdateExtraPriceInDB(updatedExtraPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return e.String(http.StatusNotFound, "extraPrice not found")
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusOK, "extraPrice updated successfully")
}
