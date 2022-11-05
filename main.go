package main

import (
	"github.com/ToggyO/dbshaker/internal"
	"log"
	"path/filepath"

	"github.com/ToggyO/dbshaker/pkg"
	_ "github.com/ToggyO/dbshaker/tests/migrations"
)

const connectionString = "host=localhost port=15436 user=dbshaker_root password=p@ssw0rd dbname=dbshaker sslmode=disable"

func main() {
	internal.Logger.Println("AHAH")

	db, err := dbshaker.OpenDbWithDriver("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	dir, err := filepath.Abs("./tests/migrations")
	if err != nil {
		log.Fatalln(err)
	}

	err = dbshaker.Down(db, dir)
	if err != nil {
		log.Fatalln(err)
	}
}
