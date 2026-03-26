package note

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	pb "notey/pkg/gen/note/v1"
	"notey/pkg/middleware"
)

type Server struct {
	pb.UnimplementedNoteServiceServer
	db  *sql.DB
	mux *http.ServeMux
}

type ErrorResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
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

	s.mux.HandleFunc("GET /notes", middleware.Logger(s.getAllNotesHandler))
}
