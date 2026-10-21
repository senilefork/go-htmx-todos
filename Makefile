TEMPL_CMD := $(HOME)/go/bin/templ

templ:
	$(TEMPL_CMD) generate

run: templ
	go run .

build: templ
	go build -o ./bin/todos .

test:
	go test ./...

vet:
	go vet ./...

include .env

MIGRATE_CMD := $(HOME)/go/bin/migrate
JET_CMD := $(HOME)/go/bin/jet
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_EXTERNAL_PORT)/$(DB_NAME)?sslmode=disable

migrate:
	$(MIGRATE_CMD) -database $(DB_URL) -path migrations up

migrate-down:
	$(MIGRATE_CMD) -database $(DB_URL) -path migrations down

migrate-down-one:
	$(MIGRATE_CMD) -database $(DB_URL) -path migrations down 1

migrate-version:
	$(MIGRATE_CMD) -database $(DB_URL) -path migrations version

migrate-create:
	@read -p "Enter migration name: " name; \
	$(MIGRATE_CMD) create -ext sql -dir migrations -seq $$name

jet-generate:
	$(JET_CMD) -dsn=$(DB_URL) -schema=public -ignore-tables=schema_migrations \
		-path=./internal/database/gen \
		-rel-model-path=../../../../../pkg/db/public/model
