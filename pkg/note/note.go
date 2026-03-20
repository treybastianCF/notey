package note

import (
	"fmt"
	"time"
)

type Note struct {
	Id        int       `json:"id"`
	Content   string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func (n Note) String() string {
	return fmt.Sprintf("%s\n\ncreated at: %s", n.Content, n.CreatedAt)
}
