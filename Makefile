docs := swagger.yaml

all: build

deps:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install golang.org/x/tools/cmd/stringer@latest

docs: $(docs)
$(docs): deps
	go generate ./...
	swag init --output . --outputTypes yaml --dir ./,../lib/swagger,./controller,../lib/models,../lib/controller

build: docs	
	go mod tidy
	go build .

run: docs	
	go run .