up:
	docker compose up -d

down:
	docker compose down

seed:
	go run cmd/seed/main.go

# ⚠️  DANGER: wipe everything and start fresh
fresh:
	@echo "⚠️  WARNING: This will DELETE all data. Continue? [y/N]" && read ans && [ $${ans:-N} = y ]
	docker compose down -v
	docker compose up -d
	sleep 3
	go run cmd/seed/main.go