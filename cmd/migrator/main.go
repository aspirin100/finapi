package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	RetriesCount = 5
	RetriesWait  = time.Second * 1
)

func main() {
	var migrationsPath string
	var postgresDSN string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&postgresDSN, "dsn", "", "URL to postgres for up migrations")

	flag.Parse()
	validateFlags(migrationsPath, postgresDSN)

	var err error

	var db *sql.DB

	for i := 0; i < RetriesCount; {
		db, err = sql.Open("postgres", postgresDSN)
		if err != nil {
			time.Sleep(RetriesWait)
			i++
		} else {
			break
		}

		if i == RetriesCount-1 {
			log.Fatal("failed to connect", err)
		}
	}

	for i := 0; i < RetriesCount; i++ {
		err = goose.Up(db, migrationsPath)
		if err != nil {
			if errors.Is(err, goose.ErrAlreadyApplied) {
				log.Print("no new migrations")
			} else {
				log.Print("failed to migrate: ", err)
			}

			time.Sleep(RetriesWait)
			i++
		} else {
			break
		}

		if i == RetriesCount-1 {
			log.Fatal(err)
		}
	}

}

func validateFlags(migrationsPath, DSN string) {
	if migrationsPath == "" {
		panic("migrations path should be not empty")
	}

	if DSN == "" {
		panic("empty DSN flag")
	}
}
