.DEFAULT_GOAL := help
UTILS_COMMAND = docker build -q -f .docker/utils/Dockerfile .docker/utils | xargs -I % docker run --rm -v .:/src %
MIGRATE_COMMAND = docker build -q -f .docker/utils/Dockerfile .docker/utils | xargs -I % docker run --network bot --rm -v .:/src %

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: confirm
confirm:
	@echo -n 'Вы уверены? [y/N] ' && read ans && [ $${ans:-N} = y ]

## help: Вывод справки по командам
.PHONY: help
help:
	@echo 'Доступные команды:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## init: Сборка рабочего окружения
init: confirm down-clear build up migrate-up

## up: Запуск контейнеров
up:
	docker-compose up -d

## down: Остановка контейнеров
down:
	docker-compose down --remove-orphans

## down-clear: Остановка контейнеров с очисткой volume (очистка всех пользовательских данных!)
down-clear: confirm
	docker-compose down -v --remove-orphans

## build: Сборка образов
build:
	docker-compose build

## bash-%: Запуск bash внутри контейнера. Пример make bash-golang
bash-%:
	docker-compose exec $* bash

## logs-%: Просмотр логов сервиса в фоллоу режиме. Пример make log-php
logs-%:
	docker-compose logs -f $*

## restart-%: Перезапуск сервиса. Пример make restart-php
restart-%:
	docker-compose restart $*

## migrate-create: Создание миграции. Пример make migrate-create name=create_table
migrate-create:
	${UTILS_COMMAND} goose create $(name) sql

## migrate-up: Выполнение миграции
migrate-up:
	${MIGRATE_COMMAND} goose up

## run/bot: run the cmd/bot application
.PHONY: run/bot
run/bot:
	@go run ./cmd/bot

## gen-enum: Генерация enum. Необходимо указать путь path до файла с перечислением
gen-enum:
	@go tool go-enum -f $(path)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
	@echo 'Running golangci-lint'
	make lint

## lint: run golangci-lint
.PHONY: lint
lint:
	@echo 'Running golangci-lint...'
	${UTILS_COMMAND} golangci-lint run ${args}

## lint-fix: run golangci-lint --fix
.PHONY: lint-fix
lint-fix:
	make lint args=--fix
# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/bot: build the cmd/bot application
.PHONY: build/bot
build/bot:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/bot ./cmd/bot
