package main

import (
	"database/sql"

	"github.com/Harman6282/attendance-app/internal/store"
	"github.com/Harman6282/attendance-app/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config     config
	store      store.Storage
	tokenMaker *token.JWTMaker
}

type config struct {
	ADDR string
	DB   *sql.DB
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", app.health)
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", app.signUpHandler)
		r.Post("/login", app.loginHandler)
		r.Get("/me", app.meHandler)
		
	})
	
	r.Route("/class", func(r chi.Router) {
		r.Post("/", app.createClass)

	})

	return r
}
