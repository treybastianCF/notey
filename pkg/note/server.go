package note

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	db  *sql.DB
	mux *http.ServeMux
}

func (s *Server) Setup(db *sql.DB, mux *http.ServeMux) {
	s.db = db
	s.mux = mux

	s.initDb()

	// setup http handlers
	s.mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello world")
	})

	s.mux.HandleFunc("DELETE /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Deleted")
	})

	s.mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "id must be an integer")
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, "Got %d", id)
	})
}
