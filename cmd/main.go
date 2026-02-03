package main

import (
	"log"
	"net/http"
)

func main() {
	db := NewConnectionPool()

	cfg := config{
		ADDR: ":8080",
		DB:   db,
	}

	app := application{
		config: cfg,
	}

	mux := app.mount()

	log.Printf("Server started at: %v", app.config.ADDR)

	err := http.ListenAndServe(app.config.ADDR, mux)

	if err != nil {
		log.Fatal(err)
	}
}
