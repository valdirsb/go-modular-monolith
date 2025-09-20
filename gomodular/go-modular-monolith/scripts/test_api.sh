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

echo "=== Teste concluído ==="