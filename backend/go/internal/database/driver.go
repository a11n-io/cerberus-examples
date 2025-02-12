package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./dbdata/example.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, err
}
