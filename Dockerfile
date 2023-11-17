FROM golang:alpine
RUN go install github.com/swaggo/swag/v2/cmd/swag@latest
# -- setup the lib
WORKDIR /monorepo/lib
COPY --from=lib go.mod go.sum ./
RUN go mod download
COPY --from=lib . .
# -- setup the service
WORKDIR /monorepo/service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
# -- compile Swagger spec
RUN swag init --output . --outputTypes yaml --dir ./,../lib/swagger,./controller,../lib/models
# -- build the service
RUN go build -o serviceStarter . 
CMD ["./serviceStarter"]