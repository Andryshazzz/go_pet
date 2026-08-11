include .env
export

export AUTH_SERVICE_ROOT=$(shell pwd)

db-up:
	@docker compose up -d postgres 

db-down:
	@docker compose down postgres

db-cleanup:
	@read -p "Clear all environment volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres && \
		rm -rf out/pgdata && \
		echo "Done"; \
	else \
		echo "Cancelled"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing seq parameter. Example: make migrate-create seq=init"; \
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
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"
	
migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

app-run:
	@export LOGGER_FOLDER=${AUTH_SERVICE_ROOT}/out/logs && \
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