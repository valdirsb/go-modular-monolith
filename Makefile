# Makefile para Go Modular Monolith

.PHONY: help build run test docker-up docker-down docker-logs clean

# Variáveis
APP_NAME=go-modular-monolith
BINARY_NAME=app
GO_VERSION=1.21

## help: Mostra esta ajuda
help:
	@echo "Comandos disponíveis:"
	@echo ""
	@echo "  build         - Compila a aplicação"
	@echo "  run           - Executa a aplicação"
	@echo "  test          - Executa os testes"
	@echo "  docker-up     - Inicia MySQL via Docker Compose"
	@echo "  docker-down   - Para e remove containers Docker"
	@echo "  docker-logs   - Mostra logs do MySQL"
	@echo "  api-test      - Executa testes da API"
	@echo "  clean         - Remove arquivos gerados"
	@echo "  setup         - Configuração inicial completa"
	@echo "  docs          - Gera documentação dos módulos"
	@echo "  db-shell      - Conecta ao MySQL via CLI"
	@echo "  dev           - Modo desenvolvimento com hot reload"
	@echo ""

## setup: Configuração inicial completa
setup:
	@echo "🚀 Configurando ambiente..."
	@cp .env.example .env
	@echo "✅ Arquivo .env criado"
	@make docker-up
	@echo "⏳ Aguardando MySQL inicializar..."
	@sleep 10
	@make build
	@echo "🎉 Setup concluído! Execute 'make run' para iniciar."

## build: Compila a aplicação
build:
	@echo "🔨 Compilando aplicação..."
	@go mod tidy
	@go build -o $(BINARY_NAME) cmd/server/main.go
	@echo "✅ Aplicação compilada: $(BINARY_NAME)"

## run: Executa a aplicação
run:
	@echo "🚀 Iniciando aplicação..."
	@go run cmd/server/main.go

## test: Executa os testes
test:
	@echo "🧪 Executando testes..."
	@go test -v ./...

## docker-up: Inicia MySQL via Docker Compose
docker-up:
	@echo "🐳 Iniciando MySQL com Docker..."
	@docker compose up -d mysql
	@echo "✅ MySQL rodando em localhost:3306"
	@echo "📊 PhpMyAdmin disponível em http://localhost:8081"

## docker-down: Para e remove containers Docker
docker-down:
	@echo "🛑 Parando containers Docker..."
	@docker-compose down

## docker-logs: Mostra logs do MySQL
docker-logs:
	@docker-compose logs -f mysql

## api-test: Executa testes da API
api-test:
	@echo "🔍 Testando API..."
	@chmod +x scripts/test_api.sh
	@./scripts/test_api.sh

## clean: Remove arquivos gerados
clean:
	@echo "🧹 Limpando arquivos..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "✅ Limpeza concluída"

## db-shell: Conecta ao MySQL via linha de comando
db-shell:
	@echo "💾 Conectando ao MySQL..."
	@docker exec -it go-modular-mysql mysql -u root -p123456 app_db

## dev: Modo desenvolvimento com hot reload (requer air)
dev:
	@if command -v air > /dev/null; then \
		echo "🔥 Iniciando em modo desenvolvimento..."; \
		air; \
	else \
		echo "❌ Air não encontrado. Instale com: go install github.com/cosmtrek/air@latest"; \
		echo "📖 Ou execute: make run"; \
	fi

## docs: Gera documentação dos módulos
docs:
	@echo "📚 Gerando documentação..."
	@echo "📖 Documentação principal: README.md"
	@echo "🏛️ Arquitetura: ARCHITECTURE.md"
	@echo "📡 API: docs/API.md"
	@echo "💾 Database: docs/DATABASE.md" 
	@echo "🛠️ Development: docs/DEVELOPMENT.md"
	@echo "🗃️ Migrations: docs/MIGRATIONS.md"
	@echo ""
	@echo "📦 Módulos:"
	@echo "   👤 User: internal/modules/user/README.md"
	@echo "   📦 Product: internal/modules/product/README.md"
	@echo "   🛒 Order: internal/modules/order/README.md"
	@echo ""
	@echo "✅ Toda documentação está atualizada!"

## stats: Estatísticas do projeto
stats:
	@echo "📊 Estatísticas do Projeto:"
	@echo ""
	@echo "📁 Linhas de código Go:"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1
	@echo ""
	@echo "📄 Arquivos de documentação:"
	@find . -name "*.md" | wc -l | xargs -I {} echo "   {} arquivos Markdown"
	@echo ""
	@echo "🧪 Arquivos de teste:"
	@find . -name "*_test.go" | wc -l | xargs -I {} echo "   {} arquivos de teste"
	@echo ""
	@echo "📦 Módulos implementados:"
	@ls internal/modules/ | wc -l | xargs -I {} echo "   {} módulos"