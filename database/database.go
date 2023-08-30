// Package  database provides functions and structs to interact with the
// postgres database
package database

import (
	"context"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	url := os.Getenv("POSTGRES_URL")
	var err error

	DB, err = pgxpool.New(context.Background(), url)

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
