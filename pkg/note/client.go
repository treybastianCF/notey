package note

type Client struct{}

func (c *Client) getAllNotes() ([]Note, error) {
	return nil, nil
}

func (c *Client) getNoteById(id int) (Note, error) {
	return Note{}, nil
}

func (c *Client) deleteNoteById(id int)
