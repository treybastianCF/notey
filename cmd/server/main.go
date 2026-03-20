package main

import (
	"context"
	"log/slog"
	"net/http"
	"notey/pkg/note"
	"notey/pkg/sql"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	db := sql.InitDB()
	defer db.Close()

	mux := http.NewServeMux()

	note := note.Server{}
	note.Setup(db, mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		slog.Info("Server started", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start %v", err)
			os.Exit(1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-exit
	slog.Info("shutdown initiated", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("server has shutdown")
}
