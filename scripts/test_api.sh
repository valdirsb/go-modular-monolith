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

echo "=== Testando API de Pedidos ==="
echo ""

# Primeiro, criar um usuário para os pedidos
echo "17. Criando usuário para testes de pedidos..."
ORDER_TIMESTAMP=$(date +%s)
ORDER_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users/" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"customer_${ORDER_TIMESTAMP}\",
    \"email\": \"customer_${ORDER_TIMESTAMP}@example.com\", 
    \"password\": \"senha123456\"
  }")

echo "$ORDER_USER_RESPONSE" | jq .

# Extrair ID do usuário
ORDER_USER_ID=$(echo "$ORDER_USER_RESPONSE" | jq -r '.id')
echo ""

if [ "$ORDER_USER_ID" != "null" ] && [ "$ORDER_USER_ID" != "" ]; then
  echo "18. Criando pedido com produtos existentes..."
  ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/orders/" \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": \"$ORDER_USER_ID\",
      \"items\": [
        {
          \"product_id\": \"prod-001\",
          \"quantity\": 1
        },
        {
          \"product_id\": \"prod-005\",
          \"quantity\": 2
        }
      ]
    }")

  echo "$ORDER_RESPONSE" | jq .

  # Extrair ID do pedido criado
  ORDER_ID=$(echo "$ORDER_RESPONSE" | jq -r '.id')
  echo ""

  if [ "$ORDER_ID" != "null" ] && [ "$ORDER_ID" != "" ]; then
    echo "19. Buscando pedido criado (ID: $ORDER_ID)..."
    curl -s "$BASE_URL/orders/$ORDER_ID" | jq .
    echo ""

    echo "20. Listando pedidos do usuário..."
    curl -s "$BASE_URL/orders/user/$ORDER_USER_ID" | jq .
    echo ""

    echo "21. Atualizando status do pedido para 'confirmed'..."
    curl -s -X PUT "$BASE_URL/orders/$ORDER_ID/status" \
      -H "Content-Type: application/json" \
      -d '{
        "status": "confirmed"
      }' | jq .
    echo ""

    echo "22. Atualizando status do pedido para 'shipped'..."
    curl -s -X PUT "$BASE_URL/orders/$ORDER_ID/status" \
      -H "Content-Type: application/json" \
      -d '{
        "status": "shipped"
      }' | jq .
    echo ""

    echo "23. Tentando cancelar pedido já enviado (deve falhar)..."
    curl -s -X POST "$BASE_URL/orders/$ORDER_ID/cancel" | jq .
    echo ""
  fi

  echo "24. Criando outro pedido para testar cancelamento..."
  CANCEL_ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/orders/" \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": \"$ORDER_USER_ID\",
      \"items\": [
        {
          \"product_id\": \"prod-002\",
          \"quantity\": 1
        }
      ]
    }")

  CANCEL_ORDER_ID=$(echo "$CANCEL_ORDER_RESPONSE" | jq -r '.id')
  echo ""

  if [ "$CANCEL_ORDER_ID" != "null" ] && [ "$CANCEL_ORDER_ID" != "" ]; then
    echo "25. Cancelando pedido pendente (ID: $CANCEL_ORDER_ID)..."
    curl -s -X POST "$BASE_URL/orders/$CANCEL_ORDER_ID/cancel" | jq .
    echo ""

    echo "26. Verificando status do pedido cancelado..."
    curl -s "$BASE_URL/orders/$CANCEL_ORDER_ID" | jq .
    echo ""
  fi

  echo "27. Testando criação de pedido com produto inexistente (deve falhar)..."
  curl -s -X POST "$BASE_URL/orders/" \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": \"$ORDER_USER_ID\",
      \"items\": [
        {
          \"product_id\": \"prod-inexistente\",
          \"quantity\": 1
        }
      ]
    }" | jq .
  echo ""

  echo "28. Testando criação de pedido com quantidade maior que o estoque (deve falhar)..."
  curl -s -X POST "$BASE_URL/orders/" \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": \"$ORDER_USER_ID\",
      \"items\": [
        {
          \"product_id\": \"prod-001\",
          \"quantity\": 999
        }
      ]
    }" | jq .
  echo ""

  echo "29. Cleanup - Deletando usuário de teste dos pedidos..."
  curl -s -X DELETE "$BASE_URL/users/$ORDER_USER_ID"
  echo "Usuário deletado"
  echo ""
fi

echo "=== Todos os testes concluídos ==="