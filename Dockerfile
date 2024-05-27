FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go

FROM alpine:3.20

COPY --from=builder /app/api /app/api

# Create a non-root user and use it to run the application
RUN adduser -D -g '' appuser
USER appuser

CMD ["/app/api"]
