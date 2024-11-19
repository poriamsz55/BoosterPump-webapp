package database

import (
	"database/sql"
	"fmt"
	"log"

	extraprice "github.com/poriamsz55/BoosterPump-webapp/internal/models/extra_price"
)

func AddExtraPriceToDB(p *extraprice.ExtraPrice) error {
	query := `INSERT INTO ` + tableExtraPrice + ` (` + columnExtraPriceName + `, ` +
		columnExtraPriceValue + `, ` +
		columnProjectIDFK + `) 
	          VALUES (?, ?, ?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Name, p.Price, p.ProjectId)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

func GetAllExtraPricesFromDB() ([]*extraprice.ExtraPrice, error) {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return nil, err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s
        FROM %s
    `, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue,
		columnProjectIDFK,
		tableExtraPrice)

	rows, err := dbHelper.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extraPrices []*extraprice.ExtraPrice
	for rows.Next() {
		var p extraprice.ExtraPrice
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.ProjectId,
		)
		if err != nil {
			return nil, err
		}
		extraPrices = append(extraPrices, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return extraPrices, nil
}

func GetExtraPriceByIdFromDB(id int) (*extraprice.ExtraPrice, error) {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return nil, err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s
        FROM %s
        WHERE %s = ?
    `, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue,
		columnProjectIDFK,
		tableExtraPrice, columnExtraPriceID)

	var p extraprice.ExtraPrice
	err = dbHelper.db.QueryRow(query, id).Scan(
		&p.Id,
		&p.Name,
		&p.Price,
		&p.ProjectId,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func GetExtraPricesByProjectIdFromDB(prjId int) ([]*extraprice.ExtraPrice, error) {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return nil, err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        SELECT %s, %s, %s
        FROM %s
        WHERE %s = ?
    `, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue,
		tableExtraPrice, columnProjectIDFK)

	var exps []*extraprice.ExtraPrice

	rows, err := dbHelper.db.Query(query, prjId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p extraprice.ExtraPrice
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
		)
		if err != nil {
			return nil, err
		}
		exps = append(exps, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exps, nil
}

func DeleteExtraPriceFromDB(id int) error {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return err
	}
	defer dbHelper.Close()

	// Check if extraPrice exists
	checkQuery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s 
        WHERE %s = ?
    `, tableExtraPrice, columnExtraPriceID)

	var count int
	err = dbHelper.db.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	query := fmt.Sprintf(`
        DELETE FROM %s 
        WHERE %s = ?
    `, tableExtraPrice, columnExtraPriceID)

	result, err := dbHelper.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateExtraPriceInDB(updatedExtraPrice *extraprice.ExtraPrice) error {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?, %s = ?
        WHERE %s = ?
    `, tableExtraPrice,
		columnExtraPriceName, columnExtraPriceValue, columnExtraPriceID)

	result, err := dbHelper.db.Exec(query,
		updatedExtraPrice.Name,
		updatedExtraPrice.Price,
		updatedExtraPrice.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func CheckExtraPriceExists(exp extraprice.ExtraPrice) (bool, error) {
	dbHelper := GetDBHelperInstance()
	err := dbHelper.Open()
	if err != nil {
		return false, err
	}
	defer dbHelper.Close()

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ? AND %s = ?
	`, tableExtraPrice, columnExtraPriceName, columnProjectIDFK)

	var count int
	err = dbHelper.db.QueryRow(query, exp.Name, exp.ProjectId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}
