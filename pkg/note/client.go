package note

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Uri string
}

func (c *Client) GetAllNotes() ([]NoteAbbr, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", c.Uri, "notes"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&er)

		return nil, fmt.Errorf("request failed status=%d msg=%s", resp.StatusCode, er.Msg)
	}
	var notes []NoteAbbr
	err = json.NewDecoder(resp.Body).Decode(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (c *Client) GetNoteById(id int) (Note, error) {
	var n Note
	resp, err := http.Get(fmt.Sprintf("%s/%s/%d", c.Uri, "notes", id))
	if err != nil {
		return n, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&er)

		return n, fmt.Errorf("request failed status=%d msg=%s", resp.StatusCode, er.Msg)
	}
	err = json.NewDecoder(resp.Body).Decode(&n)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (c *Client) DeleteNoteById(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%d", c.Uri, "notes", id), nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&er)

		return fmt.Errorf("request failed status=%d msg=%s", resp.StatusCode, er.Msg)
	}
	return nil
}

func (c *Client) CreateNote(nn NewNote) (Note, error) {
	var n Note
	req, err := json.Marshal(nn)
	if err != nil {
		return n, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/%s", c.Uri, "notes"), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return n, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		var er ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&er)

		return n, fmt.Errorf("request failed status=%d msg=%s", resp.StatusCode, er.Msg)
	}
	err = json.NewDecoder(resp.Body).Decode(&n)
	if err != nil {
		return n, err
	}
	return n, nil
}
