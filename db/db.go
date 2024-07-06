package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

func NewDB(dbDriver string, dbSource string) *sql.DB {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		fmt.Errorf("failed to open the database connection: %w", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		fmt.Errorf("failed to ping the database: %w", err)
		return nil
	}


	return db
} 