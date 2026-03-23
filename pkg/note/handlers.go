package note

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) getAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := s.getAllNotes()
	if err != nil {
		writeErrorResponse(err, &ErrorResponse{"failed to get notes", 500}, w)
		return
	}
	writeJsonResponse(notes, 200, w)
}

func (s *Server) deleteAllNotesHandler(w http.ResponseWriter, r *http.Request) {
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
}

func (s *Server) getAllNotesByIdHandler(w http.ResponseWriter, r *http.Request) {
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
}

func (s *Server) deleteNoteByIdHandler(w http.ResponseWriter, r *http.Request) {
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
}
