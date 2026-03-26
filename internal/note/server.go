package note

import (
	"context"
	"database/sql"
	"notey/internal/note/db"
	pb "notey/pkg/gen/note/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NoteServer struct {
	pb.UnimplementedNoteServiceServer
	queries *db.Queries
}

func NewNoteServer(sql *sql.DB, grpc *grpc.Server) *NoteServer {
	s := &NoteServer{
		queries: db.New(sql),
	}

	pb.RegisterNoteServiceServer(grpc, s)
	return s
}

func noteToProto(n *db.Note) *pb.Note {
	return &pb.Note{
		Id:        int32(n.ID),
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: timestamppb.New(n.Createdat.Time),
	}
}

func noteAbbrToProto(n *db.GetNotesAbbrRow) *pb.NoteAbbr {
	return &pb.NoteAbbr{
		Id:        int32(n.ID),
		Title:     n.Title,
		CreatedAt: timestamppb.New(n.Createdat.Time),
	}
}

func (s *NoteServer) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {

	note, err := s.queries.GetNote(ctx, int64(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "database error")

	}
	return &pb.GetNoteResponse{Note: noteToProto(&note)}, nil
}

func (s *NoteServer) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	if len(req.Title) < 1 || len(req.Content) < 1 {
		return nil, status.Error(codes.InvalidArgument, "title and content must not be blank")
	}

	note, err := s.queries.CreateNote(ctx, db.CreateNoteParams{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	return &pb.CreateNoteResponse{Note: noteToProto(&note)}, nil
}

func (s *NoteServer) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	err := s.queries.DeleteNote(ctx, int64(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	return &pb.DeleteNoteResponse{}, nil
}

func (s *NoteServer) GetNotes(ctx context.Context, _ *pb.GetNotesRequest) (*pb.GetNotesResponse, error) {
	notes, err := s.queries.GetNotesAbbr(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}
	res := []*pb.NoteAbbr{}
	for _, v := range notes {
		res = append(res, noteAbbrToProto(&v))
	}

	return &pb.GetNotesResponse{Notes: res}, nil
}
