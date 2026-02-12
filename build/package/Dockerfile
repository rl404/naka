# Golang base image
FROM golang:1.25 AS go_builder
LABEL stage="naka_builder"
WORKDIR /naka
COPY . .
WORKDIR /naka/cmd/naka
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o naka

# New stage from scratch
FROM jrottenberg/ffmpeg:8-alpine
LABEL stage="naka_binary"
COPY --from=go_builder /naka/cmd/naka/naka naka
ENTRYPOINT []