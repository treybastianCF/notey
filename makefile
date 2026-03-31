.PHONY: gen clean clean-db clean-build setup migrate new-note-migration check

build:
	go build -o ./out/client ./cmd/client/main.go	
	@echo "client built ./out/client"
	go build -o ./out/server ./cmd/server/main.go
	@echo "server built ./out/server"

gen:
	@echo "codegen"
	buf generate proto
	sqlc generate

clean: clean-db clean-build	

clean-db:
	@echo "removing the database"
	rm ./notey.db

clean-build:
	@echo "removing all generated files and build artifacts..."
	rm -rf pkg/gen/*
	rm -rf out/*

setup:
	@echo "installing tools"
	@echo "dont use latest in a real world scenrio!"
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate:
	@echo "running database migrations"
	goose sqlite3 ./notey.db -dir "internal/note/db/migrations/" up

new-note-migration:
	@echo "creating new migration"
	goose --dir "internal/note/db/migrations/" create $(NAME) sql


check:
	go mod tidy
	go mod verify
	gofmt -s -l . | wc -l
	go vet ./...
	gosec --exclude-generated ./...
	govulncheck ./...
