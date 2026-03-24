FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pictionary-backend ./main.go

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /app/pictionary-backend /app/pictionary-backend

USER app
EXPOSE 8080

CMD ["./pictionary-backend"]
