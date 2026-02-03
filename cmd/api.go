package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	ADDR string
	DB *sql.DB
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health check pass 🎫"))
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", app.health)

	return r
}
