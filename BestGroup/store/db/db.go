package handlingdb

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // add this
)

func NewDB(dbString string) *sql.DB {

	// Connect to database
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
