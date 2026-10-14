include .env
export

export ORDERS_SERVICE_ROOT=$(shell pwd)

db-up:
	@docker compose up -d postgres-orders

db-down:
	@docker compose down postgres-orders

db-cleanup:
	@read -p "Clear all environment volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres-orders && \
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
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-orders:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"
	
migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

kafka-up:
	@docker compose up -d kafka kafka-ui

kafka-down:
	@docker compose down kafka kafka-ui

app-run:
	@export LOGGER_FOLDER=${ORDERS_SERVICE_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export POSTGRES_PORT=5434 && \
	export KAFKA_BROKERS=localhost:9092 && \
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