FROM golang:alpine
# -- setup the lib
WORKDIR /monorepo/lib
COPY --from=lib . .
# -- setup the service
WORKDIR /monorepo/service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -- build the service
RUN go build -o serviceStarter . 
CMD ["./serviceStarter"]