all: build

deps:
	go install golang.org/x/tools/cmd/stringer@latest

build:
	go mod download
	go build .

run:
	go mod download
	go run .	


