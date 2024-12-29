RUN_DIR = ./cmd/app

.PHONY: run install clear

run:
	go run $(RUN_DIR)/main.go

install:
	go mod download

clear:
	go mod tidy