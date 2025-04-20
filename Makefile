APP_NAME = toolbox
CMD_DIR = .

.PHONY: dev test

dev:
	go run .

test:
	go test ./...

build:
	go build -o ./bin/toolbox .