package storage

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewPostgres(dbURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}
