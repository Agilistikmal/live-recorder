FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o live-recorder ./main.go

FROM alpine:latest

# Required FFmpeg for HLS downloads
RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=builder /app/live-recorder .

RUN mkdir -p /app/tmp

VOLUME ["/app/tmp"]

ENTRYPOINT ["./live-recorder"]

