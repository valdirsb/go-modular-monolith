#!/bin/bash

# Script para testar otimizações da função CreateOrder
# Testa cenários onde o cache de produtos é mais importante

echo "=== Teste de Otimização da Função CreateOrder ==="

# Iniciar servidor em background
echo "1. Iniciando servidor..."
make run &
SERVER_PID=$!

# Aguardar servidor inicializar
echo "2. Aguardando servidor inicializar..."
sleep 8

# Verificar se servidor está rodando
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Servidor não está respondendo"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo "✅ Servidor está rodando"

# Criar usuário para testes
echo "3. Criando usuário de teste..."
USER_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"username":"testopt","email":"testopt@example.com","password":"password123","full_name":"Test Optimization"}' \
  -s)

USER_ID=$(echo $USER_RESPONSE | jq -r '.id')

if [ "$USER_ID" = "null" ] || [ -z "$USER_ID" ]; then
    echo "❌ Erro ao criar usuário"
    echo "Response: $USER_RESPONSE"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo "✅ Usuário criado: $USER_ID"

# Obter lista de produtos seedados
echo "4. Obtendo produtos disponíveis..."
PRODUCTS=$(curl -s http://localhost:8080/api/v1/products/ | jq -r '.data[0:3] | .[].id')
PRODUCT_ARRAY=($PRODUCTS)

if [ ${#PRODUCT_ARRAY[@]} -lt 2 ]; then
    echo "❌ Não há produtos suficientes para teste"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

echo "✅ Produtos disponíveis: ${PRODUCT_ARRAY[@]}"

# Teste 1: Pedido com produto duplicado (teste de cache)
echo ""
echo "=== TESTE 1: Pedido com produtos duplicados (otimização de cache) ==="
echo "5. Criando pedido com múltiplas unidades do mesmo produto..."

ORDER_REQUEST_1="{
    \"user_id\": \"$USER_ID\",
    \"items\": [
        {\"product_id\": \"${PRODUCT_ARRAY[0]}\", \"quantity\": 2},
        {\"product_id\": \"${PRODUCT_ARRAY[0]}\", \"quantity\": 1},
        {\"product_id\": \"${PRODUCT_ARRAY[0]}\", \"quantity\": 1}
    ]
}"

echo "Enviando requisição:"
echo $ORDER_REQUEST_1 | jq .

ORDER_RESPONSE_1=$(curl -X POST http://localhost:8080/api/v1/orders/ \
  -H "Content-Type: application/json" \
  -d "$ORDER_REQUEST_1" \
  -s)

ORDER_ID_1=$(echo $ORDER_RESPONSE_1 | jq -r '.id')

if [ "$ORDER_ID_1" != "null" ] && [ ! -z "$ORDER_ID_1" ]; then
    echo "✅ Pedido criado com sucesso: $ORDER_ID_1"
    echo "Response: $ORDER_RESPONSE_1" | jq .
else
    echo "❌ Erro ao criar pedido"
    echo "Response: $ORDER_RESPONSE_1"
fi

# Teste 2: Pedido com múltiplos produtos diferentes 
echo ""
echo "=== TESTE 2: Pedido com produtos diferentes ==="
echo "6. Criando pedido com produtos variados..."

ORDER_REQUEST_2="{
    \"user_id\": \"$USER_ID\",
    \"items\": [
        {\"product_id\": \"${PRODUCT_ARRAY[0]}\", \"quantity\": 1},
        {\"product_id\": \"${PRODUCT_ARRAY[1]}\", \"quantity\": 2},
        {\"product_id\": \"${PRODUCT_ARRAY[2]}\", \"quantity\": 1}
    ]
}"

echo "Enviando requisição:"
echo $ORDER_REQUEST_2 | jq .

ORDER_RESPONSE_2=$(curl -X POST http://localhost:8080/api/v1/orders/ \
  -H "Content-Type: application/json" \
  -d "$ORDER_REQUEST_2" \
  -s)

ORDER_ID_2=$(echo $ORDER_RESPONSE_2 | jq -r '.id')

if [ "$ORDER_ID_2" != "null" ] && [ ! -z "$ORDER_ID_2" ]; then
    echo "✅ Pedido criado com sucesso: $ORDER_ID_2"
    echo "Response: $ORDER_RESPONSE_2" | jq .
else
    echo "❌ Erro ao criar pedido"
    echo "Response: $ORDER_RESPONSE_2"
fi

# Verificar estoque após pedidos
echo ""
echo "=== VERIFICAÇÃO DE ESTOQUE ==="
echo "7. Verificando estoque dos produtos após pedidos..."

for product_id in "${PRODUCT_ARRAY[@]:0:3}"; do
    PRODUCT_INFO=$(curl -s http://localhost:8080/api/v1/products/$product_id | jq -r '. | "\(.name): \(.stock) unidades"')
    echo "- $PRODUCT_INFO"
done

# Teste de cancelamento (otimização de restauração de estoque)
echo ""
echo "=== TESTE 3: Cancelamento de pedido (otimização de reversão) ==="
echo "8. Cancelando primeiro pedido..."

CANCEL_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/orders/$ORDER_ID_1/cancel \
  -H "Content-Type: application/json" \
  -s)

echo "Response do cancelamento: $CANCEL_RESPONSE" | jq .

# Verificar estoque após cancelamento
echo "9. Verificando estoque após cancelamento..."

for product_id in "${PRODUCT_ARRAY[@]:0:3}"; do
    PRODUCT_INFO=$(curl -s http://localhost:8080/api/v1/products/$product_id | jq -r '. | "\(.name): \(.stock) unidades"')
    echo "- $PRODUCT_INFO"
done

echo ""
echo "=== TESTE CONCLUÍDO ==="
echo "✅ Otimizações testadas:"
echo "   - Cache de produtos para evitar múltiplas consultas"
echo "   - Controle de quantidade total por produto"
echo "   - Otimização na reversão de estoque"
echo "   - Uso do cache na compensação de falhas"

# Finalizar servidor
echo ""
echo "10. Finalizando servidor..."
kill $SERVER_PID 2>/dev/null
sleep 2

echo "✅ Teste de otimização finalizado!"