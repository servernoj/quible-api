all:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install golang.org/x/tools/cmd/stringer@latest
	go mod tidy
	swag init --output . --outputTypes yaml --dir ./,../lib/swagger,./controller,../lib/models
	go generate ./...
	go build .
