# artifact:Makefile | v1.0.0 | 2026-08-19T06:54:00Z | Topic: Delivery | JobFlow developer workflow

SHELL := /bin/sh

APP_NAME := jobflow
CMD := ./cmd/jobflow

DOCKER_IMAGE ?= jobflow:dev

GO ?= go
GOFLAGS ?=
GO_TEST_FLAGS ?=

COMPOSE := docker compose
GOOSE ?= goose

DATABASE_URL ?= postgres://jobflow:jobflow@localhost:5432/jobflow?sslmode=disable

.PHONY: help
help:
	@printf '%s\n' \
		'JobFlow' \
		'' \
		'Infrastructure:' \
		'  make up              Start PostgreSQL, Redis, Redpanda' \
		'  make down            Stop infrastructure' \
		'  make restart         Restart infrastructure' \
		'  make ps              Show infrastructure status' \
		'  make logs             Follow infrastructure logs' \
		'  make clean           Remove containers and volumes' \
		'' \
		'Database:' \
		'  make migrate          Apply migrations' \
		'  make migrate-status   Show migration status' \
		'  make migrate-down     Roll back latest migration' \
		'' \
		'Development:' \
		'  make run              Run JobFlow locally' \
		'  make build            Build binary' \
		'  make fmt              Format code' \
		'  make vet              Run go vet' \
		'  make test             Run tests' \
		'  make race             Run race detector' \
		'  make coverage         Generate coverage report' \
		'' \
		'Docker:' \
		'  make image            Build Docker image' \
		'  make image-run        Run Docker image' \
		'' \
		'Verification:' \
		'  make verify            fmt + vet + test + race'

.PHONY: up
up:
	$(COMPOSE) up -d

.PHONY: down
down:
	$(COMPOSE) down

.PHONY: restart
restart:
	$(COMPOSE) restart

.PHONY: ps
ps:
	$(COMPOSE) ps

.PHONY: logs
logs:
	$(COMPOSE) logs -f

.PHONY: clean
clean:
	$(COMPOSE) down -v --remove-orphans

.PHONY: migrate
migrate:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

.PHONY: migrate-status
migrate-status:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" status

.PHONY: migrate-down
migrate-down:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" down

.PHONY: run
run:
	$(GO) run $(GOFLAGS) $(CMD)

.PHONY: build
build:
	mkdir -p bin
	$(GO) build $(GOFLAGS) -o bin/$(APP_NAME) $(CMD)

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test $(GO_TEST_FLAGS) ./...

.PHONY: race
race:
	$(GO) test -race ./...

.PHONY: coverage
coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

.PHONY: image
image:
	docker build \
		--tag $(DOCKER_IMAGE) \
		.

.PHONY: image-run
image-run:
	docker run --rm \
		--env-file .env \
		-p 8080:8080 \
		$(DOCKER_IMAGE)

.PHONY: verify
verify:
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) test
	$(MAKE) race