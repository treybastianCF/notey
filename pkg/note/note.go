package note

import (
	"fmt"
	"time"
)

type Note struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func (n Note) String() string {
	return fmt.Sprintf("%s\n\ncreated at: %s", n.Text, n.CreatedAt)
}
