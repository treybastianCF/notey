package note

import (
	"log/slog"
)

func (s *Server) initDb() {
	slog.Debug("creating notes table if needed")
	table := `
	create table if not exists notes (
		id integer not null primary key, 
		title string not null,
		content text not null,
		createdAt datetime default current_timestamp
	);
	`

	_, err := s.db.Exec(table)
	if err != nil {
		slog.Error("failed to create notes table", slog.Any("err", err), slog.String("sql", table))
	}
	slog.Debug("notes table done")
}

func (s *Server) getAllNotes() ([]NoteAbbr, error) {
	rows, err := s.db.Query("SELECT id, title, createdAt FROM notes ORDER BY createdAt DESC")
	if err != nil {
		return nil, err
	}
	res := []NoteAbbr{}
	for rows.Next() {
		var n NoteAbbr
		if err = rows.Scan(&n.Id, &n.Title, &n.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, n)
	}

	return res, nil
}

func (s *Server) getNoteById(id int) (Note, error) {
	var n Note
	stmt, err := s.db.Prepare("SELECT id, title, content, createdAt FROM notes WHERE id = ?")
	if err != nil {
		return n, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(id).Scan(&n.Id, &n.Title, &n.Content, &n.CreatedAt); err != nil {
		return n, err
	}
	return n, nil
}

func (s *Server) deleteNoteById(id int) error {
	stmt, err := s.db.Prepare("DELETE FROM notes WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(id)

	return nil
}

func (s *Server) createNote(title string, content string) (Note, error) {
	var n Note
	stmt, err := s.db.Prepare("INSERT INTO note(title, content) VALUES(?,?) RETURNING id, title, content, createdAt")
	if err != nil {
		return n, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(content).Scan(&n.Id, &n.Title, &n.Content, &n.CreatedAt); err != nil {
		return n, err
	}

	return n, nil
}
