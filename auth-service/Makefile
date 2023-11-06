all:
	go mod tidy
	swag fmt
	swag init --output swagger --outputTypes yaml
	go build .