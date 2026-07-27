include .env
export

.PHONY: env-up env-down env-cleanup

export PROJECT_ROOT = $(shell pwd)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@printf "Удалить PostgreSQL container и volume со всеми данными? [y/N] "; \
	read answer; \
	if [ "$$answer" = "y" ] || [ "$$answer" = "Y" ]; then \
		project=$$(docker compose config --format json | sed -n 's/^[[:space:]]*"name": "\([^"]*\)",*$$/\1/p' | head -n 1); \
		docker compose rm --stop --force todoapp-postgres; \
		volume=$$(docker volume ls --quiet \
			--filter "label=com.docker.compose.project=$$project" \
			--filter "label=com.docker.compose.volume=postgres-data"); \
		if [ -n "$$volume" ]; then \
			docker volume rm "$$volume"; \
		else \
			echo "PostgreSQL volume не найден."; \
		fi; \
	else \
		echo "Очистка отменена."; \
	fi
env-port-forward:
	@docker compose up -d port-forwarder
env-port-close:
	@docker compose down port-forwarder
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Ошибка: Необходимо указать имя миграции. Используйте 'make migrate-create seq=<имя_миграции>'."; \
		exit 1; \
	fi; \

	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up
migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Ошибка: Необходимо указать действие миграции. Используйте 'make migrate-up' или 'make migrate-down'."; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		$(action)