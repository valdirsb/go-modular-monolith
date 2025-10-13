#!/bin/bash

# Script para testar os serviços gRPC

echo "🧪 Testando serviços gRPC..."

# Verifica se grpcurl está instalado
if ! command -v grpcurl &> /dev/null; then
    echo "❌ grpcurl não encontrado. Instalando..."
    echo "Para instalar o grpcurl:"
    echo "go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    echo "ou"  
    echo "brew install grpcurl (no macOS)"
    echo "apt-get install grpcurl (no Ubuntu)"
    exit 1
fi

GRPC_HOST="localhost:50051"

echo "🔍 Listando serviços disponíveis..."
grpcurl -plaintext $GRPC_HOST list

echo ""
echo "📋 Listando métodos do UserService..."
grpcurl -plaintext $GRPC_HOST list api.UserService

echo ""
echo "📋 Listando métodos do ProductService..."
grpcurl -plaintext $GRPC_HOST list api.ProductService

echo ""
echo "📋 Listando métodos do OrderService..."
grpcurl -plaintext $GRPC_HOST list api.OrderService

echo ""
echo "🏥 Testando Health Check..."
grpcurl -plaintext $GRPC_HOST grpc.health.v1.Health/Check

# Teste de criação de usuário
echo ""
echo "👤 Testando criação de usuário..."
grpcurl -plaintext -d '{
  "name": "João Silva",
  "email": "joao@teste.com",
  "password": "123456"
}' $GRPC_HOST api.UserService/CreateUser

# Teste de listagem de usuários
echo ""
echo "📄 Testando listagem de usuários..."
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' $GRPC_HOST api.UserService/ListUsers

# Teste de criação de produto
echo ""
echo "📦 Testando criação de produto..."
grpcurl -plaintext -d '{
  "name": "Notebook Dell",
  "description": "Notebook Dell Inspiron 15",
  "price": 2500.99,
  "quantity": 10
}' $GRPC_HOST api.ProductService/CreateProduct

# Teste de listagem de produtos
echo ""
echo "📄 Testando listagem de produtos..."
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' $GRPC_HOST api.ProductService/ListProducts

echo ""
echo "✅ Testes gRPC concluídos!"