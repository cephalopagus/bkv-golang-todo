include .env


env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "This will remove all data in the environment. Are you sure? (y/n): " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down && \
		sudo rm -rf out/pgdata && \
		echo "Environment cleaned up."; \
	else \
		echo "Cleanup aborted."; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Migration name is required. Use 'make migrate-create seq=your_migration_name'"; \
		exit 1; \
	fi
	@docker compose run --rm todoapp-postgres-migrate create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"


migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Migration action is required. Use 'make migrate-action action=up' or 'make migrate-action action=down'"; \
		exit 1; \
	fi

	@docker compose run --rm todoapp-postgres-migrate -path=/migrations -database \
		postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

