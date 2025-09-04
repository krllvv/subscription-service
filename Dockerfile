FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sub-service ./cmd/app


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sub-service .

EXPOSE 8080

CMD ["./sub-service"]