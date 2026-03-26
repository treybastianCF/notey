package main

import (
	"log"
	"notey/internal/cli"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// not setting up TLS for dev
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(":8080", opts)
	if err != nil {
		log.Fatal("failed to initate client", err)
	}
	client := cli.NewClient(conn)
	client.Run()
}
