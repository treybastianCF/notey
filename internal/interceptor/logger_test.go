package interceptor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"notey/internal/interceptor"
)

type testStream struct{ ctx context.Context }

func (s *testStream) Context() context.Context     { return s.ctx }
func (s *testStream) SetHeader(metadata.MD) error  { return nil }
func (s *testStream) SendHeader(metadata.MD) error { return nil }
func (s *testStream) SetTrailer(metadata.MD)       {}
func (s *testStream) SendMsg(any) error            { return nil }
func (s *testStream) RecvMsg(any) error            { return nil }

// --- RequestLogger ---

func TestRequestLogger_PassesThroughResponse(t *testing.T) {
	info := &grpc.UnaryServerInfo{FullMethod: "/note.v1.NoteService/GetNote"}
	handler := func(ctx context.Context, req any) (any, error) {
		return "response", nil
	}

	resp, err := interceptor.RequestLogger(context.Background(), "request", info, handler)

	require.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestRequestLogger_PassesThroughError(t *testing.T) {
	info := &grpc.UnaryServerInfo{FullMethod: "/note.v1.NoteService/GetNote"}
	want := status.Error(codes.NotFound, "not found")
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, want
	}

	resp, err := interceptor.RequestLogger(context.Background(), nil, info, handler)

	assert.Nil(t, resp)
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestRequestLogger_PassesThroughContext(t *testing.T) {
	type key struct{}
	ctx := context.WithValue(context.Background(), key{}, "value")
	info := &grpc.UnaryServerInfo{FullMethod: "/note.v1.NoteService/CreateNote"}

	var capturedCtx context.Context
	handler := func(ctx context.Context, req any) (any, error) {
		capturedCtx = ctx
		return nil, nil
	}

	interceptor.RequestLogger(ctx, nil, info, handler)

	assert.Equal(t, "value", capturedCtx.Value(key{}))
}

func TestStreamLogger_PassesThroughNilError(t *testing.T) {
	info := &grpc.StreamServerInfo{FullMethod: "/note.v1.NoteService/WatchNotes"}
	handler := func(srv any, stream grpc.ServerStream) error { return nil }

	err := interceptor.StreamLogger(nil, &testStream{ctx: context.Background()}, info, handler)

	require.NoError(t, err)
}

func TestStreamLogger_PassesThroughError(t *testing.T) {
	info := &grpc.StreamServerInfo{FullMethod: "/note.v1.NoteService/WatchNotes"}
	want := errors.New("stream broke")
	handler := func(srv any, stream grpc.ServerStream) error { return want }

	err := interceptor.StreamLogger(nil, &testStream{ctx: context.Background()}, info, handler)

	assert.ErrorIs(t, err, want)
}

func TestStreamLogger_PassesThroughStream(t *testing.T) {
	info := &grpc.StreamServerInfo{FullMethod: "/note.v1.NoteService/WatchNotes"}
	stream := &testStream{ctx: context.Background()}

	var capturedStream grpc.ServerStream
	handler := func(srv any, ss grpc.ServerStream) error {
		capturedStream = ss
		return nil
	}

	interceptor.StreamLogger(nil, stream, info, handler)

	assert.Equal(t, stream, capturedStream)
}
