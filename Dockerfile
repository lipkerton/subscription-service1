FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o subscription-service ./cmd/app


FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/subscription-service .

EXPOSE 8080

CMD ["./subscription-service"]