# Golang base image
FROM golang:1.17 as go_builder
LABEL stage="naka_builder"
WORKDIR /naka
COPY . .
WORKDIR /naka/cmd/naka
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -mod vendor -o naka

# New stage from scratch
FROM alpine:3.14
LABEL stage="naka_binary"
RUN apk add --no-cache tzdata
COPY --from=go_builder /naka/cmd/naka/naka /cmd/naka/naka
WORKDIR /cmd/naka