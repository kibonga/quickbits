FROM golang:1.21.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/web

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/ui /app/ui

CMD ["/app/main"]