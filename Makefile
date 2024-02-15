docs := swagger.yaml

all: build

deps:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install golang.org/x/tools/cmd/stringer@latest

docs: $(docs)
$(docs): deps
	swag init \
		--output . \
		--outputTypes yaml \
		--dir ./,../lib/swagger,./controller,../lib/models,../lib/controller

build:
	go mod download
	go build .

run:
	go mod download
	go run .	


