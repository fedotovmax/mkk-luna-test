FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o api ./cmd/api

FROM alpine:3.21.3

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

ENV APP_ENV=release

COPY --from=builder /app/api ./api

COPY migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 5555

CMD ["./api"]