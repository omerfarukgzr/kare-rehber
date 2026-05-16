.PHONY: help dev backend frontend migrate seed seed-wipe reset build clean db-up db-down db-reset

help:
	@echo "KARE-REHBER — Geliştirme Komutları"
	@echo ""
	@echo "  make db-up        Postgres container'ını başlat (Docker, port 5433)"
	@echo "  make db-down      Postgres container'ını durdur"
	@echo "  make db-reset     Postgres container'ı + volume'unu silip yeniden kur"
	@echo "  make migrate      Migration'ları uygula"
	@echo "  make seed         Admin + 36 hafta + test kullanıcıları seed et"
	@echo "  make seed-wipe    Tüm tabloları temizleyip baştan seed et"
	@echo "  make reset        db-reset + migrate + seed (sıfırdan kurulum)"
	@echo "  make backend      Backend'i çalıştır (air ile hot-reload)"
	@echo "  make frontend     Frontend'i çalıştır"
	@echo "  make build        Production build (backend binary + frontend dist)"
	@echo "  make clean        tmp / dist klasörlerini temizle"

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d postgres
	@echo "Postgres ready. Now run: make migrate && make seed"

migrate:
	cd backend && go run ./cmd/migrate up

migrate-down:
	cd backend && go run ./cmd/migrate down 1

seed:
	cd backend && go run ./cmd/seed

seed-wipe:
	cd backend && go run ./cmd/seed --wipe

reset: db-reset
	@sleep 3
	cd backend && go run ./cmd/migrate up && go run ./cmd/seed

backend:
	cd backend && air 2>/dev/null || cd backend && go run ./cmd/api

frontend:
	cd frontend && npm run dev

build:
	cd backend && go build -o bin/api ./cmd/api
	cd backend && go build -o bin/migrate ./cmd/migrate
	cd backend && go build -o bin/seed ./cmd/seed
	cd frontend && npm run build

clean:
	rm -rf backend/tmp backend/bin frontend/dist frontend/node_modules/.tmp
