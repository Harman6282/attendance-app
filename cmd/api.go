package main

import (
	"database/sql"

	"github.com/Harman6282/attendance-app/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	ADDR string
	DB   *sql.DB
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", app.health)
	r.Post("/users", app.createUserHandler)

	return r
}
