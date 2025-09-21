#!/bin/bash

# Script para testar a API com MySQL

echo "=== Testando API com MySQL ==="
echo ""

# URL base da API
BASE_URL="http://localhost:8080/api/v1"

echo "1. Verificando se a aplicação está rodando..."
curl -s "http://localhost:8080/health" | jq . || echo "Aplicação não está rodando. Execute: go run cmd/server/main.go"
echo ""

echo "2. Criando um usuário..."
TIMESTAMP=$(date +%s)
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users/" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"joao_${TIMESTAMP}\",
    \"email\": \"joao_${TIMESTAMP}@example.com\", 
    \"password\": \"senha123456\"
  }")

echo "$USER_RESPONSE" | jq .

# Extrair ID do usuário criado
USER_ID=$(echo "$USER_RESPONSE" | jq -r '.id')
echo ""

if [ "$USER_ID" != "null" ] && [ "$USER_ID" != "" ]; then
  echo "3. Buscando usuário criado (ID: $USER_ID)..."
  curl -s "$BASE_URL/users/$USER_ID" | jq .
  echo ""

  echo "4. Atualizando usuário..."
  curl -s -X PUT "$BASE_URL/users/$USER_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "joao_updated"
    }' | jq .
  echo ""

  echo "5. Validando credenciais..."
  curl -s -X POST "$BASE_URL/users/validate" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"joao_${TIMESTAMP}@example.com\",
      \"password\": \"senha123456\"
    }" | jq .
  echo ""

  echo "6. Tentativa de login com senha incorreta..."
  curl -s -X POST "$BASE_URL/users/validate" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"joao_${TIMESTAMP}@example.com\",
      \"password\": \"senha_errada\"
    }" | jq .
  echo ""

  echo "7. Deletando usuário..."
  curl -s -X DELETE "$BASE_URL/users/$USER_ID"
  echo "Usuário deletado"
  echo ""

  echo "8. Tentando buscar usuário deletado..."
  curl -s "$BASE_URL/users/$USER_ID" | jq .
  echo ""
fi

echo "=== Testando API de Produtos ==="
echo ""

echo "9. Listando produtos seedados..."
PRODUCTS_RESPONSE=$(curl -s "$BASE_URL/products/")
echo "$PRODUCTS_RESPONSE" | jq '. | length' | xargs -I {} echo "Total de produtos encontrados: {}"
echo ""

echo "10. Buscando produto específico (iPhone)..."
curl -s "$BASE_URL/products/prod-001" | jq .
echo ""

echo "11. Filtrando produtos por categoria (electronics)..."
curl -s "$BASE_URL/products/?category_id=electronics" | jq .
echo ""

echo "12. Criando novo produto..."
PRODUCT_TIMESTAMP=$(date +%s)
PRODUCT_RESPONSE=$(curl -s -X POST "$BASE_URL/products/" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Produto Teste ${PRODUCT_TIMESTAMP}\",
    \"description\": \"Produto criado pelo script de teste\",
    \"price\": 99.99,
    \"stock\": 100,
    \"category_id\": \"test\"
  }")

echo "$PRODUCT_RESPONSE" | jq .

# Extrair ID do produto criado
PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | jq -r '.id')
echo ""

if [ "$PRODUCT_ID" != "null" ] && [ "$PRODUCT_ID" != "" ]; then
  echo "13. Atualizando estoque do produto criado (ID: $PRODUCT_ID)..."
  curl -s -X PUT "$BASE_URL/products/$PRODUCT_ID/stock" \
    -H "Content-Type: application/json" \
    -d '{
      "stock": 50
    }' | jq .
  echo ""

  echo "14. Atualizando produto criado..."
  curl -s -X PUT "$BASE_URL/products/$PRODUCT_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Produto Teste Atualizado",
      "price": 149.99
    }' | jq .
  echo ""

  echo "15. Deletando produto criado..."
  curl -s -X DELETE "$BASE_URL/products/$PRODUCT_ID"
  echo "Produto deletado"
  echo ""
fi

echo "16. Testando filtros avançados..."
echo "Produtos com preço entre R$ 2000 e R$ 5000:"
curl -s "$BASE_URL/products/?min_price=2000&max_price=5000" | jq '. | length' | xargs -I {} echo "Encontrados: {} produtos"
echo ""

echo "=== Todos os testes concluídos ==="