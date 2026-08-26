include .env
export

MIGRATE_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:
	go run cmd/api/main.go

dev:
	air

tidy:
	go mod tidy

seed:
	go run cmd/seeder/main.go

migrate-up:
	migrate -path migrations -database "$(MIGRATE_DSN)" up

migrate-down:
	migrate -path migrations -database "$(MIGRATE_DSN)" down 1

migrate-fresh:
	migrate -path migrations -database "$(MIGRATE_DSN)" down -all
	migrate -path migrations -database "$(MIGRATE_DSN)" up

# Example: make migrate-create name=create_product_discounts_table
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Example: make migrate-force version=29
migrate-force:
	migrate -path migrations -database "$(MIGRATE_DSN)" force $(version)
