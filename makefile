GO := go
BIN_NAME := Totion
CMD_DIR := ./cmd/totion


build:
	@$(GO) build -o $(BIN_NAME) $(CMD_DIR)

run: build
	@./$(BIN_NAME)

windows:
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BIN_NAME).exe $(CMD_DIR)