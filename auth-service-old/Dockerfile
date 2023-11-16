FROM golang:alpine
RUN go install github.com/swaggo/swag/v2/cmd/swag@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN swag init --output swagger --outputTypes yaml
RUN go build -o service .
CMD ["./service"]  