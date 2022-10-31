package main

import (
	"log"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/ToggyO/dbshaker/pkg"
	_ "github.com/ToggyO/dbshaker/tests/migrations"
)

const connectionString = "host=localhost port=15436 user=dbshaker_root password=p@ssw0rd dbname=dbshaker sslmode=disable"

// TODO: обдумать патч версии. Возмонжно, последней версией БД стоит считать последнюю примененную миграцию

// TODO: remove package `github.com/lib/pq`
func main() {
	db, err := dbshaker.OpenDbWithDriver("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	dir, err := filepath.Abs("./tests/migrations")
	if err != nil {
		log.Fatalln(err)
	}

	err = dbshaker.Up(db, dir)
	if err != nil {
		log.Fatalln(err)
	}
}
