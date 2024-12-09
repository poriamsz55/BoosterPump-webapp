package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/poriamsz55/BoosterPump-webapp/internal/database/loader"
)

func (h *DBHelper) LoadInitDatabase() error {

	// Open both databases
	javaDB, err := sql.Open("sqlite3", "./init.db")
	if err != nil {
		log.Fatal("Error opening Java database:", err)
	}
	defer javaDB.Close()

	goDB, err := sql.Open("sqlite3", "./booster_pump.db")
	if err != nil {
		log.Fatal("Error opening Go database:", err)
	}
	defer goDB.Close()

	// Convert all tables
	if err := loader.ConvertAllTables(javaDB, goDB); err != nil {
		log.Fatal("Error during conversion:", err)
	}

	// Verify all conversions
	if err := loader.VerifyAllConversions(javaDB, goDB); err != nil {
		log.Fatal("Error during verification:", err)
	}

	fmt.Println("All tables converted and verified successfully!")

	return nil
}
