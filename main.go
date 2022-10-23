package main

import (
	dbshaker "github.com/ToggyO/dbshaker/pkg"
	"log"
)

const connectionString = "host=localhost port=15436 user=dbshaker_root password=p@ssw0rd dbname=dbshaker"

func main() {
	db, err := dbshaker.OpenDbWithDriver("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}
}
