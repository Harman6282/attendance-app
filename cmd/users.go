package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Harman6282/attendance-app/internal/store"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health check pass 🎫"))
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err := app.store.Users.Create(ctx, "harman", "harman@example.com", "password123", store.Role(store.Student))
	if err != nil {
		http.Error(w, "error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
}
