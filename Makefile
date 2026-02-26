dev:
	@APP_ENV=development go run ./cmd/main.go -c .env

test:
	go test -v ./...