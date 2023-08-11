# Golang base image
FROM golang:1.21 as go_builder
LABEL stage="naka_builder"
WORKDIR /naka
COPY . .
WORKDIR /naka/cmd/naka
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o naka

# New stage from scratch
FROM alpine:3.18
LABEL stage="naka_binary"
RUN apk add --no-cache ffmpeg
COPY --from=go_builder /naka/cmd/naka/naka naka