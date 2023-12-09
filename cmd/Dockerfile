FROM golang:alpine
RUN apk add --no-cache bash
WORKDIR /cmd
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["./migrate.sh", "up"]