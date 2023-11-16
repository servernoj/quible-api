all:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go mod tidy
	go fmt ./...
	swag fmt
	swag init --output swagger --outputTypes yaml
	go build .