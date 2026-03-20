package note

import (
	"fmt"
	"log/slog"
)

func (s *Server) initDb() {
	slog.Debug("creating notes table if needed")
	table := `
	create table if not exists notes (
		id integer not null primary key, 
		conent text not null,
		createdAt datetime default current_timestamp
	);
	`

	_, err := s.db.Exec(table)
	if err != nil {
		slog.Error("failed to create notes table", slog.Any("err", err), slog.String("sql", table))
	}
	slog.Debug("notes table done")
}

func (s *Server) getAllNotes() ([]Note, error) {
	rows, err := s.db.Query("SELECT id, content, createdAt FROM notes ORDER BY createdAt DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes %v", err)
	}
	res := []Note{}
	for rows.Next() {
		var n Note
		if err = rows.Scan(&n.Id, &n.Content, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to get all notes %v", err)
		}
		res = append(res, n)
	}

	return res, nil
}

func (s *Server) getNoteById(id int) (Note, error) {
	var n Note
	stmt, err := s.db.Prepare("SELECT id, content, createdAt FROM notes WHERE id = ?")
	if err != nil {
		return n, fmt.Errorf("failed to create prepared stmt %v", err)
	}
	defer stmt.Close()
	if err = stmt.QueryRow(id).Scan(&n.Id, &n.Content, &n.CreatedAt); err != nil {
		return n, fmt.Errorf("failed to deserialize row into struct %v", err)
	}
	return n, nil
}

func (s *Server) deleteNoteById(id int) error {
	stmt, err := s.db.Prepare("DELETE FROM notes WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to delete note %v", err)
	}
	defer stmt.Close()

	stmt.Exec(id)

	return nil
}

func (s *Server) createNote(content string) (Note, error) {
	var n Note
	stmt, err := s.db.Prepare("INSERT INTO note(content) VALUES(?)")
	if err != nil {
		return n, fmt.Errorf("failed to create note %v", err)
	}
	defer stmt.Close()

	if err = stmt.QueryRow(content).Scan(&n.Id, &n.Content, &n.CreatedAt); err != nil {
		return n, fmt.Errorf("failed to create note %v", err)
	}

	return n, nil
}
