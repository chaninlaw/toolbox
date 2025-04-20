APP_NAME = toolbox

.PHONY: dev test

dev:
	go run .

test:
	go test ./...

build:
	go build -o ./bin/$(APP_NAME) .