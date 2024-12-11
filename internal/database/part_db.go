package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func AddPartToDB(p *part.Part) error {

	// Check if part already exists
	err := CheckPartFromDB(p.Name, p.Size, p.Material, p.Brand)
	if err == nil {
		return errors.New("Part already exists")
	}

	query := `INSERT INTO ` + tableParts + ` (` + columnPartName + `, ` +
		columnPartSize + `, ` +
		columnPartMaterial + `, ` +
		columnPartBrand + `, ` + columnPartPrice + `) 
	          VALUES (?, ?, ?, ?, ?)`

	stmt, err := instance.db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Name, p.Size, p.Material, p.Brand, p.Price)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	return nil
}

func CheckPartFromDB(name, size, material, brand string) error {

	// check if the part exists
	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ? AND %s = ? AND %s = ? AND %s = ?
	`, tableParts, columnPartName, columnPartSize, columnPartMaterial, columnPartBrand)

	var count int
	err := instance.db.QueryRow(query, name, size, material, brand).Scan(
		&count,
	)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("Part already exists")
	}

	return nil
}

func GetAllPartsFromDB() ([]*part.Part, error) {
	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s, %s, %s, %s
        FROM %s
    `, columnPartID, columnPartName, columnPartSize,
		columnPartMaterial, columnPartBrand, columnPartPrice,
		columnModifiedAt,
		tableParts)

	rows, err := instance.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []*part.Part
	for rows.Next() {
		var p part.Part
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Size,
			&p.Material,
			&p.Brand,
			&p.Price,
			&p.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		parts = append(parts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}

func GetPartByIdFromDB(id int) (*part.Part, error) {
	query := fmt.Sprintf(`
        SELECT %s, %s, %s, %s, %s, %s, %s
        FROM %s
        WHERE %s = ?
    `, columnPartID, columnPartName, columnPartSize,
		columnPartMaterial, columnPartBrand, columnPartPrice,
		columnModifiedAt,
		tableParts, columnPartID)

	var p part.Part
	err := instance.db.QueryRow(query, id).Scan(
		&p.Id,
		&p.Name,
		&p.Size,
		&p.Material,
		&p.Brand,
		&p.Price,
		&p.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Check if Part exists By Name
func CheckPartByNameFromDB(name string) error {

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s 
		WHERE %s = ?
	`, tableParts, columnPartName)

	var count int
	err := instance.db.QueryRow(query, name).Scan(
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

func DeletePartFromDB(id int) error {
	// Check if part exists
	checkQuery := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s 
        WHERE %s = ?
    `, tableParts, columnPartID)

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
    `, tableParts, columnPartID)

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

	// delete from device_parts
	query = fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = ?
	`, tableDeviceParts, columnPartID)

	_, err = instance.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePartInDB(updatedPart *part.Part) error {

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = ?, %s = ?, %s = ?, %s = ?, %s = ?
        WHERE %s = ?
    `, tableParts,
		columnPartName, columnPartSize, columnPartMaterial,
		columnPartBrand, columnPartPrice, columnPartID)

	result, err := instance.db.Exec(query,
		updatedPart.Name,
		updatedPart.Size,
		updatedPart.Material,
		updatedPart.Brand,
		updatedPart.Price,
		updatedPart.Id,
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
