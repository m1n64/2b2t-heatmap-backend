# Define container names
APP_CONTAINER=backend
POSTGRES_CONTAINER=db
REDIS_CONTAINER=redis

# Define dev override file
DEV_COMPOSE=docker-compose.override.yml

# Extract PostgreSQL credentials from .env file
POSTGRES_DB=$(shell grep POSTGRES_DB .env | cut -d '=' -f2)
POSTGRES_USER=$(shell grep POSTGRES_USER .env | cut -d '=' -f2)
POSTGRES_PASSWORD=$(shell grep POSTGRES_PASSWORD .env | cut -d '=' -f2)

.PHONY: help up prod prod-build stop down restart restart-container stop-container logs bash seed psql redis

## 📜 Display all available commands
help:
	@echo ""
	@echo "🔥  Available Makefile Commands 🔥"
	@echo "--------------------------------------------------------------"
	@echo "💻  Start environment:"
	@echo "  make up             - Start the DEV environment"
	@echo "  make prod           - Start the PROD environment"
	@echo "  make prod-build     - Build and Start the PROD environment"
	@echo ""
	@echo "🛑  Stop and manage containers:"
	@echo "  make stop           - Stop all containers"
	@echo "  make down           - Remove all containers and volumes"
	@echo "  make restart        - Restart all containers"
	@echo "  make restart-container CONTAINER=<name> - Restart a specific container"
	@echo "  make stop-container CONTAINER=<name>    - Stop a specific container"
	@echo ""
	@echo "🖥  Open container shell:"
	@echo "  make bash           - Open a bash shell inside the app container"
	@echo "  make seed           - Seed the database"
	@echo ""
	@echo "📜  Logs:"
	@echo "  make logs <container> - View logs of a specific container"
	@echo ""
	@echo "🐘  PostgreSQL CLI:"
	@echo "  make psql           - Open psql shell with credentials from .env"
	@echo ""
	@echo "🔥  Redis CLI:"
	@echo "  make redis          - Open redis-cli inside the Redis container"
	@echo ""

## 💻 Start the DEV environment (with override)
up:
	docker-compose -f docker-compose.yml -f $(DEV_COMPOSE) up -d

## 💻 Start the PROD environment (without override)
prod:
	docker-compose -f docker-compose.yml up -d

## 💻 Build the PROD environment (without override)
prod-build:
	docker-compose -f docker-compose.yml up -d --build

## 🛑 Stop all running containers
stop:
	docker-compose stop

## 🗑 Remove all containers and volumes
down:
	docker-compose down -v

## 🔄 Restart all containers
restart:
	docker-compose restart

## 🔄 Restart a specific container (usage: make restart-container CONTAINER=nginx)
restart-container:
	docker-compose restart $(CONTAINER)

## 🛑 Stop a specific container (usage: make stop-container CONTAINER=postgres)
stop-container:
	docker-compose stop $(CONTAINER)

## 🖥 Open a bash shell inside the app container
bash:
	docker-compose exec $(APP_CONTAINER) bash

## 📜 View logs of a specific container (usage: make logs nginx)
logs:
	docker-compose logs -f $(filter-out $@,$(MAKECMDGOALS))

## 📜 View logs of a specific container (usage: make logs nginx)
app-logs:
	docker logs gc-app -n 25

## 🐘 Open PostgreSQL shell with credentials from .env
psql:
	docker-compose exec -e PGPASSWORD=$(POSTGRES_PASSWORD) $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

## 🔥 Open Redis CLI inside the Redis container
redis:
	docker-compose exec $(REDIS_CONTAINER) redis-cli

## 🔥 Run seeders
seed:
	docker-compose exec $(APP_CONTAINER) go run cmd/seed/main.go

## Fix for make to avoid creating unnecessary files
%:
	@: