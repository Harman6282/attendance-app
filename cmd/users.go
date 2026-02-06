package main

import (
	"context"
	"github.com/Harman6282/attendance-app/internal/store"
	"log"
	"net/http"
	"strings"
	"time"
)

type usersRequestData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, "working fine", nil)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()


	var user usersRequestData
	err := app.readJSON(w, r, &user)
	if err != nil {
		log.Print(err)
	}

	if strings.TrimSpace(user.Name) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "name is required")
		return
	} else if strings.TrimSpace(user.Email) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "email is required")
		return
	} else if strings.TrimSpace(user.Password) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "password is required")
		return
	}

	res, err := app.store.Users.Create(ctx, user.Name, user.Email, user.Password, store.Role(store.Student))

	if err != nil {
		// http.Error(w, "error creating user: "+ err, http.StatusInternalServerError)
		log.Print(err)
		return
	}

	app.writeJSON(w, http.StatusOK, "user created successfully", res)
}
