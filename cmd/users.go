package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Harman6282/attendance-app/internal/errors"
	"github.com/Harman6282/attendance-app/internal/store"
	"github.com/Harman6282/attendance-app/internal/utils"
)

type signUpReqPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginUserRes struct {
	AccessToken string     `json:"accessToken"`
	User        store.User `json:"user"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, "working fine", nil)
}

func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var user signUpReqPayload
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

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	user.Password = hashed

	res, err := app.store.Users.Create(ctx, user.Name, user.Email, user.Password, store.Role(store.Student))

	if err != nil {
		if errors.IsUniqueViolation(err) {
			app.writeJSONError(w, http.StatusConflict, "email already exists")
			return
		}

		app.writeJSONError(w, http.StatusBadRequest, "failed create user")
		return
	}

	app.writeJSON(w, http.StatusOK, "user created successfully", res)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {

	var user loginRequest
	err := app.readJSON(w, r, &user)
	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, "error in reading body")
		return
	}

	if strings.TrimSpace(user.Email) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "email is required")
	} else if strings.TrimSpace(user.Password) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "password is required")
	}

	dbUser, err := app.store.Users.GetUser(r.Context(), user.Email)
	if err != nil {
		log.Print(err)
		app.writeJSONError(w, http.StatusInternalServerError, "error getting user")
		return
	}

	err = utils.CheckPassword(user.Password, dbUser.Password)
	if err != nil {
		app.writeJSONError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	accessToken, _, err := app.tokenMaker.CreateToken(dbUser.ID, dbUser.Role, 15*time.Minute)

	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	res := loginUserRes{
		AccessToken: accessToken,
		User:        *dbUser,
	}

	app.writeJSON(w, http.StatusOK, "logged in", res)
}


func (app *application) meHandler(w http.ResponseWriter, r *http.Request) {

	type meRequest struct {
		Id string `json:"id"`
	}

	var userID meRequest

	err := app.readJSON(w, r, &userID) 

	if strings.TrimSpace(userID.Id) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "id is required")
		return 
	}



	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, "error on reading id")
		return
	}


	user, err := app.store.Users.Me(r.Context(), userID.Id)
	if err != nil {
		app.writeJSONError(w, http.StatusInternalServerError, "error getting current user")
		return
	}

	app.writeJSON(w, http.StatusOK, "user found", user)
}
