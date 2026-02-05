package main

import (
	"log"
	"net/http"

	"github.com/Harman6282/attendance-app/internal/store"
)

func main() {
	db := NewConnectionPool()

	cfg := config{
		ADDR: ":8080",
		DB:   db,
	}

	store := store.NewStorage(db)
	app := application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Printf("Server started at: %v", app.config.ADDR)

	err := http.ListenAndServe(app.config.ADDR, mux)

	if err != nil {
		log.Fatal(err)
	}
}
