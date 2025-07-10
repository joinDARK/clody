run:
	cd backend/cmd && go run main.go

dev:
	cd frontend && pnpm dev

db-up:
	docker-compose up -d db

migrate:
	go run github.com/golang-migrate/migrate/v4 -path backend/migrations -database $$DATABASE_URL up

lint:
	golangci-lint run ./...