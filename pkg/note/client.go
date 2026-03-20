package note

type Client struct{}

func (c *Client) GetAllNotes() ([]NoteAbbr, error) {
	return nil, nil
}

func (c *Client) GetNoteById(id int) (Note, error) {
	return Note{}, nil
}

func (c *Client) DeleteNoteById(id int) error {
	return nil
}

func (c *Client) CreateNote(req NewNote) (Note, error) {
	return Note{}, nil
}
