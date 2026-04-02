package note_test

import (
	"context"
	"database/sql"
	"embed"
	"net"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"notey/internal/note"
	pb "notey/pkg/gen/note/v1"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/migrations/*.sql
var migrations embed.FS

func setupClient(t *testing.T) pb.NoteServiceClient {
	t.Helper()

	sqlDB, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	goose.SetBaseFS(migrations)
	require.NoError(t, goose.SetDialect("sqlite3"))
	require.NoError(t, goose.Up(sqlDB, "db/migrations"))

	srv := grpc.NewServer()
	note.NewNoteServer(sqlDB, srv)

	lis := bufconn.Listen(1024 * 1024)
	go srv.Serve(lis)
	t.Cleanup(srv.GracefulStop)
	t.Cleanup(func() { lis.Close() })

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return pb.NewNoteServiceClient(conn)
}

func TestGetNote_Success(t *testing.T) {
	client := setupClient(t)
	ctx := context.Background()

	created, err := client.CreateNote(ctx, &pb.CreateNoteRequest{Title: "Hello", Content: "World"})
	require.NoError(t, err)

	res, err := client.GetNote(ctx, &pb.GetNoteRequest{Id: created.Note.Id})

	require.NoError(t, err)
	assert.Equal(t, "Hello", res.Note.Title)
	assert.Equal(t, "World", res.Note.Content)
	assert.Equal(t, created.Note.Id, res.Note.Id)
}

func TestGetNote_NotFound(t *testing.T) {
	client := setupClient(t)

	_, err := client.GetNote(context.Background(), &pb.GetNoteRequest{Id: 999})

	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestCreateNote_Success(t *testing.T) {
	client := setupClient(t)

	res, err := client.CreateNote(context.Background(), &pb.CreateNoteRequest{
		Title:   "My Note",
		Content: "Some content",
	})

	require.NoError(t, err)
	assert.Equal(t, "My Note", res.Note.Title)
	assert.Equal(t, "Some content", res.Note.Content)
	assert.NotZero(t, res.Note.Id)
	assert.NotNil(t, res.Note.CreatedAt)
}

func TestCreateNote_BlankTitle(t *testing.T) {
	client := setupClient(t)

	_, err := client.CreateNote(context.Background(), &pb.CreateNoteRequest{Title: "", Content: "content"})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateNote_BlankContent(t *testing.T) {
	client := setupClient(t)

	_, err := client.CreateNote(context.Background(), &pb.CreateNoteRequest{Title: "title", Content: ""})

	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateNote_BroadcastsToWatchers(t *testing.T) {
	client := setupClient(t)
	ctx := t.Context()

	stream, err := client.WatchNotes(ctx, &pb.WatchNotesRequest{})
	require.NoError(t, err)

	received := make(chan string, 1)
	go func() {
		if res, err := stream.Recv(); err == nil {
			received <- res.Note.Title
		}
	}()

	// Allow the server goroutine to register the subscriber before broadcasting.
	time.Sleep(10 * time.Millisecond)

	_, err = client.CreateNote(context.Background(), &pb.CreateNoteRequest{Title: "Live", Content: "Update"})
	require.NoError(t, err)

	select {
	case title := <-received:
		assert.Equal(t, "Live", title)
	case <-time.After(time.Second):
		t.Fatal("watcher did not receive broadcast")
	}
}

func TestDeleteNote_Success(t *testing.T) {
	client := setupClient(t)
	ctx := context.Background()

	created, err := client.CreateNote(ctx, &pb.CreateNoteRequest{Title: "To Delete", Content: "bye"})
	require.NoError(t, err)

	_, err = client.DeleteNote(ctx, &pb.DeleteNoteRequest{Id: created.Note.Id})
	require.NoError(t, err)

	_, err = client.GetNote(ctx, &pb.GetNoteRequest{Id: created.Note.Id})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestGetNotes_ReturnsAllNotes(t *testing.T) {
	client := setupClient(t)
	ctx := context.Background()

	for _, title := range []string{"First", "Second", "Third"} {
		_, err := client.CreateNote(ctx, &pb.CreateNoteRequest{Title: title, Content: "body"})
		require.NoError(t, err)
	}

	res, err := client.GetNotes(ctx, &pb.GetNotesRequest{})

	require.NoError(t, err)
	assert.Len(t, res.Notes, 3)
}

func TestGetNotes_EmptyDB(t *testing.T) {
	client := setupClient(t)

	res, err := client.GetNotes(context.Background(), &pb.GetNotesRequest{})

	require.NoError(t, err)
	assert.Empty(t, res.Notes)
}

func TestWatchNotes_ExitsOnContextCancel(t *testing.T) {
	client := setupClient(t)
	ctx, cancel := context.WithCancel(context.Background())

	stream, err := client.WatchNotes(ctx, &pb.WatchNotesRequest{})
	require.NoError(t, err)

	cancel()

	_, err = stream.Recv()
	require.Error(t, err)
	assert.Equal(t, codes.Canceled, status.Code(err))
}

func TestWatchNotes_ReceivesMultipleNotes(t *testing.T) {
	client := setupClient(t)
	ctx := t.Context()

	stream, err := client.WatchNotes(ctx, &pb.WatchNotesRequest{})
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	titles := []string{"A", "B", "C"}
	for _, title := range titles {
		_, err := client.CreateNote(context.Background(), &pb.CreateNoteRequest{Title: title, Content: "body"})
		require.NoError(t, err)
	}

	for _, expected := range titles {
		res, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, expected, res.Note.Title)
	}
}
