# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_ADDRESS ?= 127.0.0.1:5432
POSTGRES_DATABASE ?= article

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)

# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile
include ./misc/make/help.Makefile

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: dev-env dev-air             ## Startup / Spinup Docker Compose and air
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

install-deps: goose air gotestsum tparse mockery ## Install Development Dependencies (locally).
deps: $(GOOSE) $(AIR) $(GOTESTSUM) $(TPARSE) $(MOCKERY) $(GOLANGCI) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev-env: ## Bootstrap Environment (with Docker-Compose help).
	@ docker-compose up -d --build postgres

dev-env-test: dev-env ## Run application (with Docker-Compose help)
	@ $(MAKE) image-build
	docker-compose up web

dev-air: $(AIR) ## Starts AIR (Continuous Development app).
	air

docker-stop:
	@ docker-compose down

docker-teardown:
	@ docker-compose down --remove-orphans -v

# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

build: ## Builds binary
	@ printf "Building application... "
	@ go build \
		-trimpath  \
		-o engine \
		./app/
	@ echo "done"

build-race: ## Builds binary (with -race flag)
	@ printf "Building application with race flag... "
	@ go build \
		-trimpath  \
		-race      \
		-o engine \
		./app/
	@ echo "done"

go-generate: $(MOCKERY) ## Runs go generate ./...
	go generate ./...

TESTS_ARGS := --format testname --jsonfile gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile   coverage.out
TESTS_ARGS += -test.timeout        5s
TESTS_ARGS += -race

tests: $(GOTESTSUM)
	@ gotestsum $(TESTS_ARGS) -short

tests-complete: tests $(TPARSE) ## Run Tests & parse details
	@cat gotestsum.json.out | $(TPARSE) -all -notests

# ~~~ Docker Build ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.ONESHELL:
image-build:
	@ echo "Docker Build"
	@ DOCKER_BUILDKIT=0 docker build \
		--file Dockerfile \
		--tag otus-hla-homework \
			.

# ~~~ Database Migrations using Goose ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

GOOSE := $(shell which goose)
MIGRATIONS_PATH := $(PWD)/migrations
POSTGRES_DSN := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable"

migrate-up: $(GOOSE) ## Apply all migrations up.
	@ $(GOOSE) -dir $(MIGRATIONS_PATH) postgres $(POSTGRES_DSN) up

.PHONY: migrate-down
migrate-down: $(GOOSE) ## Apply all migrations down.
	@ $(GOOSE) -dir $(MIGRATIONS_PATH) postgres $(POSTGRES_DSN) down

.PHONY: migrate-status
migrate-status: $(GOOSE) ## Display current migration status.
	@ $(GOOSE) -dir $(MIGRATIONS_PATH) postgres $(POSTGRES_DSN) status

.PHONY: migrate-create
migrate-create: $(GOOSE) ## Create a new migration with a specified name.
	@ read -p "Please provide name for the migration: " Name; \
	$(GOOSE) -dir $(MIGRATIONS_PATH) create "$${Name}" sql

.PHONY: migrate-reset
migrate-reset: $(GOOSE) ## Resets the database (drops everything and migrates up again).
	@ $(GOOSE) -dir $(MIGRATIONS_PATH) postgres $(POSTGRES_DSN) reset

# ~~~ Cleans ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleaning artifacts... "
	@rm -f *.out
	@echo "done."

clean-docker: ## Removes dangling docker images
	@ docker image prune -f