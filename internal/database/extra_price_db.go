package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	extraprice "github.com/poriamsz55/BoosterPump-webapp/internal/models/extra_price"
	tehrantime "github.com/poriamsz55/BoosterPump-webapp/internal/time"
)

func AddExtraPriceToDB(prjId int, expName string, expValue uint64) error {

	// check if the extra price already exists
	err := extraPriceExists(prjId, expName)
	if err != nil {
		return errors.New("extra price already exists")

	}

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

	_, err = stmt.Exec(expName, expValue, prjId)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

func extraPriceExists(prjId int, expName string) error {

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE %s = ? AND %s = ?
	`, tableExtraPrice, columnProjectIDFK, columnExtraPriceName)

	var count int
	err := instance.db.QueryRow(query, prjId, expName).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("extra price already exists")
	}

	return nil
}

func GetAllExtraPricesFromDB() ([]*extraprice.ExtraPrice, error) {

	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s, %s
        FROM %s
    `, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue,
		columnProjectIDFK, columnModifiedAt,
		tableExtraPrice)

	rows, err := instance.db.Query(query)
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
			&p.ModifiedAt,
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

	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s, %s
        FROM %s
        WHERE %s = ?
    `, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue,
		columnProjectIDFK, columnModifiedAt,
		tableExtraPrice, columnExtraPriceID)

	var p extraprice.ExtraPrice
	err := instance.db.QueryRow(query, id).Scan(
		&p.Id,
		&p.Name,
		&p.Price,
		&p.ProjectId,
		&p.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}

	p.ModifiedAt = tehrantime.FormattedDateTime(p.ModifiedAt)
	return &p, nil
}

func GetExtraPricesByProjectIdFromDB(db *sql.DB, prjId int) ([]*extraprice.ExtraPrice, error) {

	if db == nil {
		db = instance.db
	}

	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s
        FROM %s
        WHERE %s = ?
    `, columnExtraPriceID, columnExtraPriceName, columnExtraPriceValue, columnModifiedAt,
		tableExtraPrice, columnProjectIDFK)

	var exps []*extraprice.ExtraPrice

	rows, err := db.Query(query, prjId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p extraprice.ExtraPrice
		p.ProjectId = prjId
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		p.ModifiedAt = tehrantime.FormattedDateTime(p.ModifiedAt)
		exps = append(exps, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exps, nil
}

// CheckExtraPriceByNameFromDB
func CheckExtraPriceByNameFromDB(expName string) error {
	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ?
	`, tableExtraPrice, columnExtraPriceName)

	var count int
	err := instance.db.QueryRow(query, expName).Scan(
		&count,
	)
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteExtraPriceFromDB(id int) error {

	// Check if extraPrice exists
	checkQuery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s 
        WHERE %s = ?
    `, tableExtraPrice, columnExtraPriceID)

	var count int
	err := instance.db.QueryRow(checkQuery, id).Scan(&count)
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

	result, err := instance.db.Exec(query, id)
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

// DeleteExtraPricesByProjectId
func DeleteExtraPricesByProjectId(prjId int) error {

	query := fmt.Sprintf(`
		DELETE FROM %s 
		WHERE %s = ?
	`, tableExtraPrice, columnProjectIDFK)

	_, err := instance.db.Exec(query, prjId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateExtraPriceInDB(updatedExtraPrice *extraprice.ExtraPrice) error {

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?, %s = ?
        WHERE %s = ?
    `, tableExtraPrice,
		columnExtraPriceName, columnExtraPriceValue, columnExtraPriceID)

	result, err := instance.db.Exec(query,
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

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ? AND %s = ?
	`, tableExtraPrice, columnExtraPriceName, columnProjectIDFK)

	var count int
	err := instance.db.QueryRow(query, exp.Name, exp.ProjectId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}
