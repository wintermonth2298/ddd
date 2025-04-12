package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func mustMigrateUp(db *sqlx.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose set dialect: %v", err)
	}

	if err := goose.Up(db.DB, "./internal/catalog/migrations"); err != nil {
		log.Fatalf("goose up: %v", err)
	}
}
