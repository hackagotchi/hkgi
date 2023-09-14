// Package  database provides functions and structs to interact with the
// postgres database
package database

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	url := os.Getenv("POSTGRES_URL")
	var err error

	DB, err = sqlx.Connect("postgres", url)

	if err != nil {
		log.Fatal(err)
	}

	runMigrations()

	log.Println("Successfully connected to database")
}

func runMigrations() {
	m, err := migrate.New("file://./migrations", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	m.Steps(2)
}
