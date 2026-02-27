VERSION=

STEPS ?= 0

dev:
	@APP_ENV=development go run ./cmd/api/main.go -c .env

test:
	go test -v ./...

migrator-up:
	@APP_ENV=development go run ./cmd/migrator/main.go -m up -steps $(STEPS) -c .env

migrator-down:
	@APP_ENV=development go run ./cmd/migrator/main.go -m down -steps $(STEPS) -c .env

migrator-force:
	@APP_ENV=development go run ./cmd/migrator/main.go -m force -version $(VERSION) -c .env

migrator-version:
	@APP_ENV=development go run ./cmd/migrator/main.go -m version -c .env

swagger:
	swag init -g cmd/api/main.go

# docker exec -it mariadb bash
# mariadb -u mkk_luna_owner -p mkk_luna