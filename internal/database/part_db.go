package database

import (
	"log"

	"github.com/poriamsz55/BoosterPump-webapp/internal/models/part"
)

func AddPartToDB(p *part.Part) error {
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
