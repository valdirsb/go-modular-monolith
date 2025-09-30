# Demo do Framework - Sistema de Pagamentos

Vamos criar um exemplo prático usando o template de pagamentos:

```bash
# 1. Criar projeto de pagamentos
./scripts/create_project.sh sistema-pagamentos payment

# 2. Entrar no diretório
cd sistema-pagamentos

# 3. Configurar ambiente
make setup-framework
make docker-up

# 4. Executar
make run-framework
```

## Endpoints do Sistema

### Framework Endpoints
- `GET /health` - Health check
- `GET /framework/info` - Informações do framework

### Business Endpoints (se módulos ativos)
- `POST /api/v1/users/` - Criar usuário
- `GET /api/v1/users/:id` - Buscar usuário  
- `POST /api/v1/orders/` - Criar pedido
- `GET /api/v1/orders/:id` - Buscar pedido

## Exemplo de Resposta do Framework

### GET /framework/info
```json
{
  "framework": {
    "name": "sistema-pagamentos", 
    "version": "1.0.0"
  },
  "features": {
    "core": ["http_server", "event_system", "logging"],
    "database": ["mysql", "redis"],
    "integrations": ["stripe", "pix", "mercado_pago"],
    "modules": ["user", "billing", "order"],
    "security": ["cors", "rate_limiting", "encryption"]
  }
}
```

### GET /health
```json
{
  "status": "healthy",
  "services": {
    "framework": null,
    "mysql": null
  },
  "features": {
    "integrations": ["stripe", "pix"],
    "modules": ["user", "order"]
  }
}
```

## Como Usar para Novos Projetos

1. **Configure** framework.yaml com suas necessidades
2. **Execute** make run-framework 
3. **Desenvolva** módulos específicos
4. **Deploy** apenas com funcionalidades ativas
5. **Monitore** via health checks

O framework está pronto para ser sua base de desenvolvimento!