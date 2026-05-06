include .env
export

export PROJECT_ROOT=$(shell pwd)


env-up:
	@docker compose up -d postgres 

env-down:
	@docker compose down postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres && \
		rm -rf out/pgdata $$ \
		echo "Готово"; \
	else \
		echo "Отменено"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down -d port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi;

	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5433/${POSTGRES_DB}?sslmode=disable \
		"$(action)"
	
migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export POSTGRES_PORT=5433 && \
	go mod tidy && \
	go run cmd/main.go

app-deploy:
	@docker compose up -d --build deploy-app

app-undeploy:
	@docker compose down deploy-app

ps:
	@docker compose ps

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/main.go \
		-o docs \
		--parseInternal \
		--parseDependency