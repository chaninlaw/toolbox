APP_NAME = toolbox
CMD_DIR = ./cmd/$(APP_NAME)

.PHONY: dev test

dev:
	go run ./cmd/toolbox/

test:
	go test ./...
