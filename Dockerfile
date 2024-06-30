# Build stage
FROM golang:1.18-alpine AS builder
LABEL maintainer="doziestar"
WORKDIR /app
COPY . .
RUN go build -o datavinci ./cmd/datavinci

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/datavinci .
EXPOSE 8080
CMD ["./datavinci"]