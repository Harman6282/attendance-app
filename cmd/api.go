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
	})

	r.Group(func(r chi.Router) {
		r.Use(app.authMiddleware)
		r.Get("/users/me", app.meHandler)
		r.Route("/class", func(r chi.Router) {
			r.With(app.teacherOnly).Post("/", app.createClass)
			r.With(app.requireRole(store.Teacher, store.Student)).Get("/:id", app.getClass)
			r.With(app.teacherOnly).Patch("/", app.addStudent)
			r.With(app.requireRole(store.Teacher, store.Student)).Get("/:id/my-attendance", app.myAttendance)
		})

		r.With(app.teacherOnly).Post("/attendance/start", app.startAttendance)

		r.Route("/students", func(r chi.Router) {
			r.With(app.teacherOnly).Get("/", app.getAllStudents)
		})
	})

	return r
}
