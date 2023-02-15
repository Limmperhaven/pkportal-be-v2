# Параметры подключения к БД из конфига
DBUSER := $(shell yq '.psql.user' etc/config.yml)
DBPASSWORD := $(shell yq '.psql.pass' etc/config.yml)
DBHOST := $(shell yq '.psql.host' etc/config.yml)
DBPORT := $(shell yq '.psql.port' etc/config.yml)
DBNAME := $(shell yq '.psql.dbname' etc/config.yml)
SSLMODE := $(shell yq '.psql.sslmode' etc/config.yml)

default: help

# Зависимости
dep:
	go mod tidy
	go mod vendor
.PHONY: dep

# FMT & GOIMPORT
fmt:
	go fmt ./... && goimports -w .
.PHONY: fmt

prepare:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/mikefarah/yq/v4@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
.PHONY: prepare

# Миграции вперед до конца
mig-up:
	migrate -source file://migrations/postgres -database 'postgres://$(DBUSER):$(DBPASSWORD)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=$(SSLMODE)' up
.PHONY: mig-up

mig-dn:
	migrate -source file://migrations/postgres -database 'postgres://$(DBUSER):$(DBPASSWORD)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=$(SSLMODE)' down 1
.PHONY: mig-dn

mig-reset:
	migrate -source file://migrations/postgres -database 'postgres://$(DBUSER):$(DBPASSWORD)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=$(SSLMODE)' drop
.PHONY: mig-reset

sql-generate:
	go run internal/models/generate.go --step1
	sqlboiler -d --wipe --add-enum-types --no-tests -o internal/models/tpportal -p tpportal -c etc/config.yml psql
	go run internal/models/generate.go --step2
.PHONY: sql-generate

test-up:
	docker compose -p pkportal-test -f docker-compose.yaml up -d --wait
.PHONY: test-up

test-dn:
	docker compose -p pkportal-test -f docker-compose.yaml stop
	docker compose -p pkportal-test -f docker-compose.yaml rm -f
.PHONY: test-dn

env-up:
	docker compose -p pkportal -f docker-compose-dev.yaml up -d --wait
.PHONY: env-up

env-dn:
	docker compose -p pkportal -f docker-compose-dev.yaml stop
	docker compose -p pkportal -f docker-compose-dev.yaml rm -f
.PHONY: env-dn

h:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    dep			- Исправление зависимостей"
	@echo "    fmt			- Форматирование кодовой базы"
	@echo "    prepare		- Инициализация утилиты для миграций и sqlboiler"
	@echo "    mig-reset	- Миграция назад до конца"
	@echo "    mig-up		- Миграция вперёд до конца"
	@echo "    mig-dn		- Миграция назад на одну"
	@echo "    sql-generate - Генерация сущностей для БД"

.PHONY: h
help: h
.PHONY: help