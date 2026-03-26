package note

import (
	"context"
	"database/sql"
	"notey/internal/note/db"
	pb "notey/pkg/gen/note/v1"
)

type NoteServer struct {
	pb.UnimplementedNoteServiceServer
	queries *db.Queries
}

func NewNoteServer(sql *sql.DB) *NoteServer {
	return &NoteServer{
		queries: db.New(sql),
	}
}

func (s *NoteServer) GetNote(ctx context.Context, req *pb.GetNoteRequest) {

}
