# -------------------------------
# 🛠️ Project Metadata
# -------------------------------
MODULE := $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "v0.0.1")
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

SHELL := /bin/bash

# -------------------------------
# 📁 Config and DB Connection
# -------------------------------
CONFIG_FILE ?= ./config/config.yaml

DB_HOST := $(shell sed -n 's/^ *host:[[:space:]]*//p' $(CONFIG_FILE) | head -n 1 | tr -d '[:space:]')
DB_PORT := $(shell sed -n 's/^ *port:[[:space:]]*//p' $(CONFIG_FILE) | head -n 1 | tr -d '[:space:]')
DB_USER := $(shell sed -n 's/^ *user:[[:space:]]*//p' $(CONFIG_FILE) | head -n 1 | tr -d '[:space:]')
DB_PASSWORD := $(shell sed -n 's/^ *password:[[:space:]]*//p' $(CONFIG_FILE) | head -n 1 | tr -d '[:space:]')
DB_NAME := $(shell sed -n 's/^ *db_name:[[:space:]]*//p' $(CONFIG_FILE) | head -n 1 | tr -d '[:space:]')


# APP_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
APP_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# -------------------------------
# 📦 Migrations
# -------------------------------
MIGRATIONS_DIR := shared/clients/db/migrations
MIGRATE := migrate -path $(MIGRATIONS_DIR) -database $(APP_DSN)

# -------------------------------
# 🔁 Migration Commands
# -------------------------------
.PHONY: migrate
migrate: ## Run all migrations
	@echo "🚀 Running migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	@echo "↩️ Rolling back last migration..."
	@$(MIGRATE) down 1

.PHONY: migrate-reset
migrate-reset: ## Drop all migrations and re-run
	@echo "🧨 Dropping all tables..."
	@$(MIGRATE) drop -f
	@echo "🚀 Re-running migrations..."
	@$(MIGRATE) up

.PHONY: migrate-new
migrate-new: ## Create a new migration file
	@if [ -z "$(name)" ]; then \
		read -p "Enter migration name: " name; \
	else \
		name="$(name)"; \
	fi; \
	name_clean=$$(echo $$name | tr ' ' '_' | tr A-Z a-z); \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $$name_clean
.PHONY: print-dsn
print-dsn:
	@echo $(APP_DSN)
