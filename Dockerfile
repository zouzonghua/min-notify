FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o min-notify main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/min-notify .
COPY static ./static
VOLUME ["/app/config.json"]
EXPOSE 5001
ENTRYPOINT ["./min-notify"]
