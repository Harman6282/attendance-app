package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (app *application) getClass(w http.ResponseWriter, r *http.Request) {
	classId := chi.URLParam(r, "id")

	class, err := app.store.Classes.Get(r.Context(), classId)
	if err != nil {
		log.Print(err)
		app.writeJSONError(w, http.StatusInternalServerError, "error fetching class")
		return
	}

	app.writeJSON(w, http.StatusOK, "class fetched successfully", class)
}

func (app *application) myAttendance(w http.ResponseWriter, r *http.Request) {
	classId := chi.URLParam(r, "id")
	if strings.TrimSpace(classId) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "class id is required")
		return
	}

	type myAttendanceRequest struct {
		StudentID string `json:"student_id"`
	}

	studentID := strings.TrimSpace(r.URL.Query().Get("student_id"))
	
	if studentID == "" {
		studentID = strings.TrimSpace(app.userIDFromAuthorization(r))
	}

	if studentID == "" {
		var req myAttendanceRequest
		err := app.readJSON(w, r, &req)
		if err != nil && !errors.Is(err, io.EOF) {
			app.writeJSONError(w, http.StatusBadRequest, "failed to read json from body")
			return
		}
		studentID = strings.TrimSpace(req.StudentID)
	}

	if studentID == "" {
		app.writeJSONError(w, http.StatusBadRequest, "student id is required")
		return
	}

	attendance, err := app.store.Classes.GetMyAttendance(r.Context(), classId, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.writeJSONError(w, http.StatusNotFound, "attendance not found")
			return
		}
		log.Print(err)
		app.writeJSONError(w, http.StatusInternalServerError, "error fetching attendance")
		return
	}

	app.writeJSON(w, http.StatusOK, "my attendance", attendance)
}

type startAttendanceRequest struct {
	ClassID string `json:"class_id"`
}

func (app *application) startAttendance(w http.ResponseWriter, r *http.Request) {
	var req startAttendanceRequest

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ClassID) == "" {
		app.writeJSONError(w, http.StatusBadRequest, "classId is required")
		return
	}

	session, err := app.store.Classes.StartAttendance(r.Context(), req.ClassID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.writeJSONError(w, http.StatusNotFound, "class not found")
			return
		}
		log.Print(err)
		app.writeJSONError(w, http.StatusInternalServerError, "error starting attendance")
		return
	}
	app.writeJSON(w, http.StatusOK, "attendance start", session)
}

func (app *application) userIDFromAuthorization(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	claims, err := app.tokenMaker.VerifyToken(strings.TrimSpace(parts[1]))
	if err != nil {
		return ""
	}

	return claims.ID
}
