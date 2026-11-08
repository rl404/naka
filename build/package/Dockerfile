# Golang base image
FROM golang:1.27.3 AS go_builder
LABEL stage="naka_builder"
WORKDIR /naka
COPY . .
WORKDIR /naka/cmd/naka
RUN CGO_ENABLED=0 go build -mod=vendor -ldflags="-s -w" -o naka

# New stage from scratch
FROM jrottenberg/ffmpeg:8-alpine
LABEL stage="naka_binary"
RUN adduser -D -u 10001 naka
COPY --from=go_builder /naka/cmd/naka/naka /usr/local/bin/naka
WORKDIR /home/naka
USER naka
ENTRYPOINT ["/usr/local/bin/naka"]
