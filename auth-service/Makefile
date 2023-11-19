all:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go mod tidy
	swag init --output . --outputTypes yaml --dir ./,../lib/swagger,./controller,../lib/models
	go build .
