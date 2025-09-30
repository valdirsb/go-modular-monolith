#!/bin/bash

# Script para criar novos projetos baseados no Go Modular Framework
# Usage: ./scripts/create_project.sh <project-name> [template]

set -e

PROJECT_NAME="$1"
TEMPLATE="${2:-basic}"

if [ -z "$PROJECT_NAME" ]; then
    echo "❌ Nome do projeto é obrigatório"
    echo "Usage: $0 <project-name> [template]"
    echo ""
    echo "Templates disponíveis:"
    echo "  basic    - Projeto básico com HTTP + MySQL"
    echo "  api      - API REST completa"
    echo "  grpc     - Serviço gRPC"
    echo "  payment  - Sistema de pagamentos"
    echo "  maps     - Integração com mapas"
    echo "  full     - Todas as funcionalidades ativadas"
    exit 1
fi

echo "🚀 Criando projeto: $PROJECT_NAME"
echo "📋 Template: $template"

# Criar diretório do projeto
mkdir -p "$PROJECT_NAME"
cd "$PROJECT_NAME"

# Copiar estrutura base
cp -r ../cmd .
cp -r ../internal .
cp -r ../pkg .
cp ../go.mod.template go.mod
cp ../go.sum .
cp ../Makefile .
cp ../docker-compose.yml .
cp ../.env.example .env

# Atualizar go.mod com o nome do projeto
sed -i "s/go-modular-monolith/$PROJECT_NAME/g" go.mod

# Gerar framework.yaml baseado no template
case "$TEMPLATE" in
    "basic")
        generate_basic_config
        ;;
    "api")
        generate_api_config
        ;;
    "grpc")
        generate_grpc_config
        ;;
    "payment")
        generate_payment_config
        ;;
    "maps")
        generate_maps_config
        ;;
    "full")
        generate_full_config
        ;;
    *)
        echo "❌ Template '$TEMPLATE' não encontrado"
        exit 1
        ;;
esac

# Atualizar imports nos arquivos Go
find . -name "*.go" -exec sed -i "s/go-modular-monolith/$PROJECT_NAME/g" {} \;

echo "✅ Projeto criado com sucesso!"
echo ""
echo "Próximos passos:"
echo "1. cd $PROJECT_NAME"
echo "2. make setup-framework"
echo "3. make docker-up"
echo "4. make run-framework"
echo ""

generate_basic_config() {
    cat > framework.yaml << 'EOF'
framework:
  name: "projeto-basico"
  version: "1.0.0"

core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true
  database_migrations: true

database:
  mysql: true
  postgresql: false
  mongodb: false
  sqlite: false

cache:
  redis: false
  memcached: false
  in_memory: false

protocols:
  http: true
  grpc: false
  websockets: false
  graphql: false

integrations:
  jwt: true
  oauth2: false
  firebase_auth: false
  stripe: false
  google_maps: false
  
security:
  cors: true
  rate_limiting: false

development:
  swagger: false
  debug_mode: true

modules:
  user: true
  product: true
  order: false
EOF
}

generate_api_config() {
    cat > framework.yaml << 'EOF'
framework:
  name: "api-rest"
  version: "1.0.0"

core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true
  database_migrations: true

database:
  mysql: true
  postgresql: false

cache:
  redis: true
  in_memory: false

protocols:
  http: true
  grpc: false

integrations:
  jwt: true
  oauth2: true
  
security:
  cors: true
  rate_limiting: true
  csrf_protection: true

development:
  swagger: true
  debug_mode: true

modules:
  user: true
  product: true
  order: true
  notification: true
EOF
}

generate_payment_config() {
    cat > framework.yaml << 'EOF'
framework:
  name: "sistema-pagamentos"
  version: "1.0.0"

core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true
  database_migrations: true

database:
  mysql: true

cache:
  redis: true

protocols:
  http: true

integrations:
  jwt: true
  stripe: true
  paypal: true
  mercado_pago: true
  pix: true
  
security:
  cors: true
  rate_limiting: true
  encryption: true

modules:
  user: true
  billing: true
  order: true
EOF
}

generate_maps_config() {
    cat > framework.yaml << 'EOF'
framework:
  name: "sistema-mapas"
  version: "1.0.0"

core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true

database:
  mysql: true

integrations:
  jwt: true
  google_maps: true
  mapbox: false
  
modules:
  user: true
  location: true
EOF
}

generate_full_config() {
    cat > framework.yaml << 'EOF'
framework:
  name: "projeto-completo"
  version: "1.0.0"

core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true
  database_migrations: true

database:
  mysql: true
  postgresql: true

cache:
  redis: true
  memcached: true

protocols:
  http: true
  grpc: true
  websockets: true

integrations:
  jwt: true
  oauth2: true
  firebase_auth: true
  stripe: true
  google_maps: true
  rabbitmq: true
  sendgrid: true
  prometheus: true
  
security:
  cors: true
  rate_limiting: true
  encryption: true

development:
  swagger: true
  debug_mode: true

modules:
  user: true
  product: true
  order: true
  notification: true
  analytics: true
  audit: true
EOF
}