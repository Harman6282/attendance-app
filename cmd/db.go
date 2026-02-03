package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const DB_URL = "postgres://admin:adminpassword@localhost:5432/attendance-db?sslmode=disable"

func NewConnectionPool() *sql.DB {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("db ping failed")
	}

	log.Print("Connection pool established")
	return db
}
