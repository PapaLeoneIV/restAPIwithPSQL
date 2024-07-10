package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
) 



func NewDB(dbDriver string, dbSource string) *sql.DB {
	fmt.Printf("Opening the database\n")
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Printf("failed to open the database connection: %v\n", err)
		return nil
	}
	fmt.Printf("Pinging Database\n")
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Printf("failed to ping the database: %v\n", err)
		return nil
	}
	fmt.Printf("Database Connected\n")

	return db
} 
