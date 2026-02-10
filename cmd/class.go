package main

import (
	"log"
	"net/http"
	"strings"
)

func (app *application) createClass(w http.ResponseWriter, r *http.Request) {

	type createClass struct {
		ClassName string `json:"class_name"`
		TeacherId string `json:"teacher_id"`
	}

	var createInput createClass

	err := app.readJSON(w, r, &createInput)
	if err != nil {
		log.Print(err)
		app.writeJSONError(w, http.StatusBadRequest, "error reading class name")
		return
	}

	if strings.TrimSpace(createInput.ClassName) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "class name is required")
		return
	}
	if strings.TrimSpace(createInput.TeacherId) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "teacher id is required")
		return
	}

	res, err := app.store.Classes.Create(r.Context(), createInput.ClassName, createInput.TeacherId)

	if err != nil {
		log.Print(err)
		app.writeJSONError(w, http.StatusInternalServerError, "error creating class")
		return
	}

	app.writeJSON(w, http.StatusOK, "class created successfully", res)

}
