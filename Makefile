# Makefile для ЗооПлатформы
# Использование: make <target>
# Справка: make help

.PHONY: help main admin petbase shelter owner volunteer clinic all stop clean install

# Цвета для вывода
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
RED    := \033[0;31m
NC     := \033[0m

# По умолчанию показываем help
.DEFAULT_GOAL := help

##@ Основные команды

help: ## Показать эту справку
	@echo "$(BLUE)🐾 ЗооПлатформа - Команды разработки$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Использование:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BLUE)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Запуск сервисов

main: ## Запустить Main (соцсеть) - порты 3000/8000
	@echo "$(GREEN)🚀 Starting Main...$(NC)"
	@echo "$(YELLOW)   Frontend: http://localhost:3000$(NC)"
	@echo "$(YELLOW)   Backend:  http://localhost:8000$(NC)"
	@cd main/backend && air & \
	cd main/frontend && npm run dev

admin: ## Запустить Admin (Main + Admin) - порты 4000/9000
	@echo "$(GREEN)🚀 Starting Main + Admin...$(NC)"
	@echo "$(YELLOW)   Main:  http://localhost:3000 / :8000$(NC)"
	@echo "$(YELLOW)   Admin: http://localhost:4000 / :9000$(NC)"
	@cd main/backend && air & \
	cd admin/backend && air & \
	cd admin/frontend && npm run dev

petbase: ## Запустить PetBase (Main + PetBase) - порты 4100/8100
	@echo "$(GREEN)🚀 Starting Main + PetBase...$(NC)"
	@echo "$(YELLOW)   Main:    http://localhost:3000 / :8000$(NC)"
	@echo "$(YELLOW)   PetBase: http://localhost:4100 / :8100$(NC)"
	@cd main/backend && air & \
	cd petbase/backend && air & \
	cd petbase/frontend && npm run dev

shelter: ## Запустить Shelter (Main + Admin + Shelter) - порты 5100/8200
	@echo "$(GREEN)🚀 Starting Main + Admin + Shelter...$(NC)"
	@echo "$(YELLOW)   Main:    http://localhost:3000 / :8000$(NC)"
	@echo "$(YELLOW)   Admin:   http://localhost:4000 / :9000$(NC)"
	@echo "$(YELLOW)   Shelter: http://localhost:5100 / :8200$(NC)"
	@cd main/backend && air & \
	cd admin/backend && air & \
	cd shelter/backend && air & \
	cd shelter/frontend && npm run dev

owner: ## Запустить Owner (Main + Owner) - порты 6100/8400
	@echo "$(GREEN)🚀 Starting Main + Owner...$(NC)"
	@echo "$(YELLOW)   Main:  http://localhost:3000 / :8000$(NC)"
	@echo "$(YELLOW)   Owner: http://localhost:6100 / :8400$(NC)"
	@cd main/backend && air & \
	cd owner/backend && air & \
	cd owner/frontend && npm run dev

volunteer: ## Запустить Volunteer (Main + Volunteer) - порты 6200/8500
	@echo "$(GREEN)🚀 Starting Main + Volunteer...$(NC)"
	@echo "$(YELLOW)   Main:      http://localhost:3000 / :8000$(NC)"
	@echo "$(YELLOW)   Volunteer: http://localhost:6200 / :8500$(NC)"
	@cd main/backend && air & \
	cd volunteer/backend && air & \
	cd volunteer/frontend && npm run dev

clinic: ## Запустить Clinic (Main + Admin + Clinic) - порты 6300/8600
	@echo "$(GREEN)🚀 Starting Main + Admin + Clinic...$(NC)"
	@echo "$(YELLOW)   Main:   http://localhost:3000 / :8000$(NC)"
	@echo "$(YELLOW)   Admin:  http://localhost:4000 / :9000$(NC)"
	@echo "$(YELLOW)   Clinic: http://localhost:6300 / :8600$(NC)"
	@cd main/backend && air & \
	cd admin/backend && air & \
	cd clinic/backend && air & \
	cd clinic/frontend && npm run dev

all: ## Запустить ВСЕ сервисы (как ./run)
	@echo "$(GREEN)🚀 Starting ALL services...$(NC)"
	@./run

##@ Управление процессами

stop: ## Остановить все процессы
	@echo "$(YELLOW)🛑 Stopping all services...$(NC)"
	@pkill -f "air" 2>/dev/null || true
	@pkill -f "next dev" 2>/dev/null || true
	@pkill -f "npm run dev" 2>/dev/null || true
	@echo "$(GREEN)✅ All services stopped$(NC)"

restart-main: stop main ## Перезапустить Main

restart-admin: stop admin ## Перезапустить Admin

restart-clinic: stop clinic ## Перезапустить Clinic

restart-owner: stop owner ## Перезапустить Owner

restart-volunteer: stop volunteer ## Перезапустить Volunteer

##@ Установка и настройка

install: ## Установить все npm зависимости
	@echo "$(GREEN)📦 Installing npm dependencies...$(NC)"
	@cd main/frontend && npm install && \
	cd ../../admin/frontend && npm install && \
	cd ../../petbase/frontend && npm install && \
	cd ../../shelter/frontend && npm install && \
	cd ../../owner/frontend && npm install && \
	cd ../../volunteer/frontend && npm install && \
	cd ../../clinic/frontend && npm install && \
	cd ../../shared && npm install
	@echo "$(GREEN)✅ All npm dependencies installed$(NC)"

install-go: ## Установить все Go зависимости
	@echo "$(GREEN)📦 Installing Go dependencies...$(NC)"
	@cd main/backend && go mod download && \
	cd ../../admin/backend && go mod download && \
	cd ../../petbase/backend && go mod download && \
	cd ../../shelter/backend && go mod download && \
	cd ../../owner/backend && go mod download && \
	cd ../../volunteer/backend && go mod download && \
	cd ../../clinic/backend && go mod download && \
	cd ../../database && go mod download
	@echo "$(GREEN)✅ Go dependencies installed$(NC)"

install-all: install install-go ## Установить все зависимости (npm + Go)

##@ База данных

db-backup: ## Создать резервную копию БД
	@echo "$(GREEN)💾 Creating database backup...$(NC)"
	@./scripts/backup-database.sh

db-restore: ## Восстановить БД из последнего бэкапа
	@echo "$(YELLOW)⚠️  Restoring database from backup...$(NC)"
	@./scripts/restore-database.sh

db-migrate: ## Применить миграции
	@echo "$(GREEN)🔄 Running database migrations...$(NC)"
	@cd database && go run migrate.go

db-status: ## Показать статус БД
	@echo "$(BLUE)📊 Database status:$(NC)"
	@ls -lh database/data.db 2>/dev/null || echo "$(RED)Database not found$(NC)"
	@echo ""
	@echo "$(BLUE)Recent backups:$(NC)"
	@ls -lht database/backups/ 2>/dev/null | head -5 || echo "$(YELLOW)No backups found$(NC)"

##@ Тестирование

test: ## Запустить все тесты
	@echo "$(GREEN)🧪 Running all tests...$(NC)"
	@./tests/performance-test.sh

test-api: ## Тестировать API
	@echo "$(GREEN)🧪 Testing API endpoints...$(NC)"
	@cd tests && ./api-test.sh

test-sso: ## Тестировать SSO
	@echo "$(GREEN)🧪 Testing SSO integration...$(NC)"
	@./tests/sso-test.sh

test-performance: ## Тест производительности
	@echo "$(GREEN)🧪 Running performance tests...$(NC)"
	@./tests/performance-test.sh

##@ Очистка

clean: ## Очистить кэш и временные файлы
	@echo "$(YELLOW)🧹 Cleaning cache and temporary files...$(NC)"
	@find . -name ".next" -type d -prune -exec rm -rf {} + 2>/dev/null || true
	@find . -name "tmp" -type d -path "*/backend/tmp" -prune -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)✅ Cleaned$(NC)"

clean-all: clean ## Очистить всё (включая node_modules)
	@echo "$(RED)🧹 Cleaning everything (including node_modules)...$(NC)"
	@echo "$(YELLOW)This may take a while...$(NC)"
	@find . -name "node_modules" -type d -prune -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)✅ Everything cleaned. Run 'make install' to reinstall$(NC)"

clean-cache: ## Очистить только кэш Next.js
	@echo "$(YELLOW)🧹 Cleaning Next.js cache...$(NC)"
	@find . -name ".next" -type d -prune -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)✅ Cache cleaned$(NC)"

clean-logs: ## Очистить логи
	@echo "$(YELLOW)🧹 Cleaning logs...$(NC)"
	@find logs -name "*.log" -type f -delete 2>/dev/null || true
	@echo "$(GREEN)✅ Logs cleaned$(NC)"

##@ Docker

docker-build: ## Собрать Docker образы
	@echo "$(GREEN)🐳 Building Docker images...$(NC)"
	@cd infrastructure && docker-compose build

docker-up: ## Запустить в Docker
	@echo "$(GREEN)🐳 Starting Docker containers...$(NC)"
	@cd infrastructure && docker-compose up -d

docker-down: ## Остановить Docker
	@echo "$(YELLOW)🐳 Stopping Docker containers...$(NC)"
	@cd infrastructure && docker-compose down

docker-logs: ## Показать логи Docker
	@cd infrastructure && docker-compose logs -f

docker-restart: docker-down docker-up ## Перезапустить Docker

##@ Разработка

dev-main: ## Режим разработки Main (с логами)
	@echo "$(GREEN)🔧 Development mode: Main$(NC)"
	@mkdir -p logs/main
	@cd main/backend && air 2>&1 | tee ../../logs/main/backend.log & \
	cd main/frontend && npm run dev 2>&1 | tee ../../logs/main/frontend.log

dev-admin: ## Режим разработки Admin (с логами)
	@echo "$(GREEN)🔧 Development mode: Admin$(NC)"
	@mkdir -p logs/admin logs/main
	@cd main/backend && air 2>&1 | tee ../../logs/main/backend.log & \
	cd admin/backend && air 2>&1 | tee ../../logs/admin/backend.log & \
	cd admin/frontend && npm run dev 2>&1 | tee ../../logs/admin/frontend.log

lint-go: ## Проверить Go код
	@echo "$(GREEN)🔍 Linting Go code...$(NC)"
	@cd main/backend && golangci-lint run || echo "$(YELLOW)golangci-lint not installed$(NC)"

lint-ts: ## Проверить TypeScript код
	@echo "$(GREEN)🔍 Linting TypeScript code...$(NC)"
	@cd main/frontend && npm run lint

format-go: ## Форматировать Go код
	@echo "$(GREEN)✨ Formatting Go code...$(NC)"
	@find . -name "*.go" -not -path "*/vendor/*" -not -path "*/node_modules/*" -exec gofmt -w {} \;
	@echo "$(GREEN)✅ Go code formatted$(NC)"

format-ts: ## Форматировать TypeScript код
	@echo "$(GREEN)✨ Formatting TypeScript code...$(NC)"
	@cd main/frontend && npm run format || echo "$(YELLOW)Format script not found$(NC)"

##@ Git

git-status: ## Показать статус Git
	@git status

git-log: ## Показать последние 10 коммитов
	@git log --oneline --decorate -10

git-tags: ## Показать все теги
	@git tag -l

git-pull: ## Обновить из репозитория
	@echo "$(GREEN)📥 Pulling from repository...$(NC)"
	@git pull

git-push: ## Отправить в репозиторий (с тегами)
	@echo "$(GREEN)📤 Pushing to repository...$(NC)"
	@git push && git push --tags

##@ Информация

ports: ## Показать все порты сервисов
	@echo "$(BLUE)📊 Порты сервисов:$(NC)"
	@echo "  $(GREEN)Main:$(NC)      3000 (frontend) / 8000 (backend)"
	@echo "  $(GREEN)Admin:$(NC)     4000 (frontend) / 9000 (backend)"
	@echo "  $(GREEN)PetBase:$(NC)   4100 (frontend) / 8100 (backend)"
	@echo "  $(GREEN)Shelter:$(NC)   5100 (frontend) / 8200 (backend)"
	@echo "  $(GREEN)Owner:$(NC)     6100 (frontend) / 8400 (backend)"
	@echo "  $(GREEN)Volunteer:$(NC) 6200 (frontend) / 8500 (backend)"
	@echo "  $(GREEN)Clinic:$(NC)    6300 (frontend) / 8600 (backend)"

status: ## Показать статус всех процессов
	@echo "$(BLUE)📊 Статус процессов:$(NC)"
	@ps aux | grep -E "(air|next dev)" | grep -v grep || echo "  $(YELLOW)Нет запущенных процессов$(NC)"

check-ports: ## Проверить занятые порты
	@echo "$(BLUE)📊 Занятые порты:$(NC)"
	@lsof -i :3000 -i :4000 -i :4100 -i :5100 -i :6100 -i :6200 -i :6300 -i :8000 -i :8100 -i :8200 -i :8400 -i :8500 -i :8600 -i :9000 2>/dev/null || echo "  $(GREEN)Все порты свободны$(NC)"

version: ## Показать версию проекта
	@echo "$(BLUE)🐾 ЗооПлатформа$(NC)"
	@echo "$(GREEN)Version: v0.8.0$(NC)"
	@git describe --tags --always 2>/dev/null || echo "No git tags"

info: ## Показать информацию о проекте
	@echo "$(BLUE)🐾 ЗооПлатформа - Информация$(NC)"
	@echo ""
	@echo "$(GREEN)Версия:$(NC) v0.8.0"
	@echo "$(GREEN)Микросервисы:$(NC) 7 (Main, Admin, PetBase, Shelter, Owner, Volunteer, Clinic)"
	@echo "$(GREEN)База данных:$(NC) SQLite (database/data.db)"
	@echo "$(GREEN)Документация:$(NC) docs/"
	@echo ""
	@echo "$(YELLOW)Быстрый старт:$(NC)"
	@echo "  make install      # Установить зависимости"
	@echo "  make admin        # Запустить Admin"
	@echo "  make stop         # Остановить всё"
	@echo ""
	@echo "$(YELLOW)Полная справка:$(NC) make help"

##@ Утилиты

logs-main: ## Показать логи Main
	@tail -f logs/main/*.log 2>/dev/null || echo "$(YELLOW)Логи не найдены. Запустите 'make dev-main'$(NC)"

logs-admin: ## Показать логи Admin
	@tail -f logs/admin/*.log 2>/dev/null || echo "$(YELLOW)Логи не найдены. Запустите 'make dev-admin'$(NC)"

backup-all: ## Полный бэкап (БД + uploads)
	@echo "$(GREEN)💾 Creating full backup...$(NC)"
	@./scripts/backup-all.sh

health: ## Проверить здоровье всех сервисов
	@echo "$(BLUE)🏥 Checking services health...$(NC)"
	@curl -s http://localhost:8000/api/health 2>/dev/null && echo "$(GREEN)✅ Main API$(NC)" || echo "$(RED)❌ Main API$(NC)"
	@curl -s http://localhost:9000/api/health 2>/dev/null && echo "$(GREEN)✅ Admin API$(NC)" || echo "$(RED)❌ Admin API$(NC)"
	@curl -s http://localhost:8100/api/health 2>/dev/null && echo "$(GREEN)✅ PetBase API$(NC)" || echo "$(RED)❌ PetBase API$(NC)"
	@curl -s http://localhost:8200/api/health 2>/dev/null && echo "$(GREEN)✅ Shelter API$(NC)" || echo "$(RED)❌ Shelter API$(NC)"
	@curl -s http://localhost:8400/api/health 2>/dev/null && echo "$(GREEN)✅ Owner API$(NC)" || echo "$(RED)❌ Owner API$(NC)"
	@curl -s http://localhost:8500/api/health 2>/dev/null && echo "$(GREEN)✅ Volunteer API$(NC)" || echo "$(RED)❌ Volunteer API$(NC)"
	@curl -s http://localhost:8600/api/health 2>/dev/null && echo "$(GREEN)✅ Clinic API$(NC)" || echo "$(RED)❌ Clinic API$(NC)"
