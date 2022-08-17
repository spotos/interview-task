.PHONY: init start stop ssh logs

init:
	./scripts/docker/init.sh

start: ## Start Docker containers
	docker compose up -d

stop: ## Stop Docker containers
	docker compose stop

ssh: ## SSH into api container
	docker compose run --rm api sh

logs: ## Show container logs
	docker compose logs -f api

help: ## Display available commands
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
