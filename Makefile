.PHONY: help up down logs ps generate
.DEFAULT_GOAL := help

up: ## Do docker compose up with hot reload
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

generate: ## Generate codes
	go generate ./...

migrate: ## Run migration
	mysqldef -u gpt -p gpt -h 127.0.0.1 -P 33306 gpt < ./_tools/mysql/schema.sql

dry-migrate: ## Run migration
	mysqldef -u gpt -p gpt -h 127.0.0.1 -P 33306 gpt --dry-run < ./_tools/mysql/schema.sql

help: ## Show options
	@grep -E '^[a-zA-Z_]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
