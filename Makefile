# Makefile Lapanganku — perintah pengembangan lokal.
# Jalankan `make` untuk melihat daftar perintah.

SHELL := /bin/bash
MIGRATE_DEV_DB := postgres://lapanganku:lapanganku@localhost:5432/lapanganku?sslmode=disable

.PHONY: help db db-down db-reset migrate backend frontend check build

help: ## Tampilkan daftar perintah
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

db: ## Start PostgreSQL dev + jalankan migration
	docker compose -f compose.dev.yaml up -d

db-down: ## Hentikan PostgreSQL dev
	docker compose -f compose.dev.yaml down

db-reset: ## Hapus data dev lalu start ulang + migration
	docker compose -f compose.dev.yaml down -v
	docker compose -f compose.dev.yaml up -d

migrate: ## Jalankan migration dev (manual)
	docker compose -f compose.dev.yaml run --rm migrate

backend: ## Jalankan backend Go (baca .env.local)
	cd backend && set -a && source ../.env.local && set +a && go run .

frontend: ## Jalankan frontend SvelteKit dev server
	cd frontend && bun run dev --host 0.0.0.0

check: ## Verifikasi: backend test + frontend check
	cd backend && go test ./...
	cd frontend && bun run check

build: ## Build produksi frontend
	cd frontend && bun run build
