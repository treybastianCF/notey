package main

import (
	"log/slog"
	"net"
	"notey/internal/interceptor"
	"notey/internal/note"
	"notey/internal/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {

	db := sql.InitDB()
	defer db.Close()

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.RequestLogger),
		grpc.StreamInterceptor(interceptor.StreamLogger))

	// don't actually need to access it once it's init
	_ = note.NewNoteServer(db, grpcSrv)

	s, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("failed to start server", slog.Any("err", err))
		os.Exit(1)
	}

	go func() {
		slog.Info("Server started", slog.String("addr", ":8080"))
		if err := grpcSrv.Serve(s); err != nil {
			slog.Error("server failed to start", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-exit
	slog.Info("shutdown initiated", slog.String("signal", sig.String()))

	stopped := make(chan struct{})
	go func() {
		grpcSrv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		slog.Info("server has shutdown gracefully")
	case <-time.After(5 * time.Second):
		slog.Warn("shutdown timedout, forcing shutdown NOW")
		grpcSrv.Stop()
	}

}
