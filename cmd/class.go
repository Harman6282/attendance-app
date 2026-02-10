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

func (app *application) addStudent(w http.ResponseWriter, r *http.Request) {

	type addStudentRequest struct {
		ClassId   string `json:"class_id"`
		StudentId string `json:"student_id"`
	}

	var addStudent addStudentRequest

	err := app.readJSON(w, r, &addStudent)
	if err != nil {
		app.writeJSONError(w, http.StatusBadRequest, "error reading add student json")
		return
	}

	if strings.TrimSpace(addStudent.ClassId) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "class Id is required")
		return
	}

	if strings.TrimSpace(addStudent.StudentId) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "student Id is required")
		return
	}

	res, err := app.store.Classes.AddStudent(r.Context(), addStudent.StudentId, addStudent.ClassId)
	if err != nil {
		log.Print(err)
		app.writeJSONError(w, http.StatusInternalServerError, "error while adding student in class")
		return
	}

	app.writeJSON(w, http.StatusOK, "student added", res)
}
