FROM golang:1.17-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o demo .

FROM alpine

COPY --from=builder ["/build/demo", "/"]

ENTRYPOINT ["/demo"]