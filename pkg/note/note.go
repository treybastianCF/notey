package note

import (
	"fmt"
	"time"
)

type Note struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func (n *Note) String() string {
	return fmt.Sprintf("%s\n%s\n\ncreated at: %s\n", n.Title, n.Content, n.CreatedAt)
}

type NoteAbbr struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

func (n *NoteAbbr) String() string {
	return fmt.Sprintf("%d\t%s\t%s\n", n.Id, n.CreatedAt, n.Title)
}
