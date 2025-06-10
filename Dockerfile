## Build
FROM golang:1.23-bookworm AS build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download -x
COPY . ./

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o application

## Run
FROM alpine:latest

WORKDIR /app
COPY --from=build /build/application ./

EXPOSE 8888

ENTRYPOINT ["./application"]
