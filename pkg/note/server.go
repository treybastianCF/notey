package note

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"notey/pkg/middleware"
	"strconv"
)

type Server struct {
	db  *sql.DB
	mux *http.ServeMux
}

type ErrorResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

type NewNote struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func writeErrorResponse(err error, er *ErrorResponse, w http.ResponseWriter) {
	slog.Error("returning error response", slog.Any("err", err))
	resp, err := json.Marshal(er)
	if err != nil {
		w.WriteHeader(500)
		slog.Error("failed to seralize error response", slog.Any("err", err))
		return
	}
	w.WriteHeader(er.Status)
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(500)
		slog.Error("failed to write error response", slog.Any("err", err))
	}
}

func writeJsonResponse(obj any, code int, w http.ResponseWriter) {
	resp, err := json.Marshal(obj)
	if err != nil {
		writeErrorResponse(err, &ErrorResponse{"failed to serialize response", 500}, w)
		return
	}
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		writeErrorResponse(err, &ErrorResponse{"failed to write response", 500}, w)
		return
	}
}

func (s *Server) Setup(db *sql.DB, mux *http.ServeMux) {
	s.db = db
	s.mux = mux

	s.initDb()

	// setup http handlers TODO: refactor these to seprate functions / file
	s.mux.HandleFunc("GET /n,otes", middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		notes, err := s.getAllNotes()
		if err != nil {
			writeErrorResponse(err, &ErrorResponse{"failed to get notes", 500}, w)
			return
		}
		writeJsonResponse(notes, 200, w)
	}))

	s.mux.HandleFunc("DELETE /notes/{id}", middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeErrorResponse(err, &ErrorResponse{"not found", 404}, w)
			return
		}
		if err = s.deleteNoteById(id); err != nil {
			writeErrorResponse(err, &ErrorResponse{"failed to delete note", 500}, w)
			return
		}
		w.WriteHeader(204)
	}))

	s.mux.HandleFunc("GET /notes/{id}", middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeErrorResponse(err, &ErrorResponse{"not found", 404}, w)
			return
		}
		note, err := s.getNoteById(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeErrorResponse(err, &ErrorResponse{"not found", 404}, w)
				return
			}

			slog.Info("err", slog.Any("err", err))
			writeErrorResponse(err, &ErrorResponse{"failed to get note", 500}, w)
			return
		}

		writeJsonResponse(note, 200, w)
	}))

	s.mux.HandleFunc("POST /notes", middleware.Logger(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req NewNote
		if err := decoder.Decode(&req); err != nil {
			writeErrorResponse(err, &ErrorResponse{"failed to process request", 400}, w)
			return
		}
		// light validation
		if len(req.Content) < 1 || len(req.Title) < 1 {
			writeErrorResponse(fmt.Errorf("validation failure"), &ErrorResponse{"content cannot be blank", 400}, w)
			return
		}
		note, err := s.createNote(req.Title, req.Content)
		if err != nil {
			writeErrorResponse(err, &ErrorResponse{"failed to create note", 500}, w)
			return
		}
		writeJsonResponse(note, 201, w)
	}))
}
