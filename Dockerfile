FROM golang:1.23.7-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=local
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 7000
CMD ["./main"]
