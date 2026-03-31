package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func RequestLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)
	st, _ := status.FromError(err)
	slog.Info("requeset", slog.String("method", info.FullMethod),
		slog.String("duration", time.Since(start).String()),
		slog.String("code", st.Code().String()), slog.Any("error", err))
	return resp, err
}

func StreamLogger(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	slog.Info("rpc stream started", slog.String("method", info.FullMethod))
	err := handler(srv, ss)

	st, _ := status.FromError(err)

	slog.Info("rpc stream ended",
		slog.String("method", info.FullMethod),
		slog.String("durration", time.Since(start).String()),
		slog.String("code", st.Code().String()),
		slog.Any("error", err))

	return err
}
