FROM golang:alpine
RUN go install github.com/swaggo/swag/v2/cmd/swag@latest
RUN go install golang.org/x/tools/cmd/stringer@latest
# -- setup the lib
WORKDIR /monorepo/lib
COPY --from=lib . .
# -- setup the service
WORKDIR /monorepo/service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -- compile Swagger spec
RUN swag init \
  --output . \
  --outputTypes yaml \
  --dir ./,../lib/swagger,./controller,../lib/models,../lib/controller
# -- generate code tagged with //go:generate
RUN go generate ./...
# -- build the service
RUN go build -o serviceStarter . 
CMD ["./serviceStarter"]