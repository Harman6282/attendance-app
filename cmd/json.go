package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    *any   `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (app *application) writeJSONError(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var res ErrorResponse
	res.Success = false
	res.Message = message

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return fmt.Errorf("error while encoding error response: %w", err)

	}
	return nil

}

func (app *application) writeJSON(w http.ResponseWriter, status int, message string, data any) (*SuccessResponse, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var res SuccessResponse
	res.Success = true
	res.Message = message
	res.Data = &data

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxbytes := 1_048_576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxbytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(dst)

}
