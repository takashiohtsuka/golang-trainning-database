build:
	go build -o bin/migrate ./cmd/migrate/main.go

migrate:
	go run ./cmd/migrate/main.go

seed:
	go run ./seed/ -mode=init

seed-add-women:
	go run ./seed/ -mode=add-women

seed-large:
	go run ./seed/large/

.PHONY: build migrate seed seed-add-women seed-large
