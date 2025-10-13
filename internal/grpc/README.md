# gRPC Implementation

Este projeto implementa serviços gRPC para todos os módulos seguindo a arquitetura modular estabelecida.

## 📋 Serviços Disponíveis

### UserService (api.UserService)
- `CreateUser` - Cria um novo usuário
- `GetUser` - Obtém usuário por ID
- `UpdateUser` - Atualiza um usuário
- `DeleteUser` - Remove um usuário
- `ListUsers` - Lista usuários com paginação
- `GetUserByEmail` - Obtém usuário por email

### ProductService (api.ProductService)
- `CreateProduct` - Cria um novo produto
- `GetProduct` - Obtém produto por ID
- `UpdateProduct` - Atualiza um produto
- `DeleteProduct` - Remove um produto
- `ListProducts` - Lista produtos com paginação
- `SearchProducts` - Busca produtos por critério

### OrderService (api.OrderService)
- `CreateOrder` - Cria um novo pedido
- `GetOrder` - Obtém pedido por ID
- `UpdateOrderStatus` - Atualiza status do pedido
- `DeleteOrder` - Cancela um pedido
- `ListOrders` - Lista pedidos com paginação
- `GetOrdersByUser` - Obtém pedidos por usuário

## 🚀 Como Executar

1. **Inicie o banco de dados:**
   ```bash
   make docker-up
   ```

2. **Execute a aplicação com framework (inclui gRPC):**
   ```bash
   make run-framework
   ```

3. **Teste os serviços gRPC:**
   ```bash
   make grpc-test
   ```

## 🔧 Configuração

### Portas
- **HTTP**: 8080 (configurável via `PORT` no .env)
- **gRPC**: 50051 (configurável via `GRPC_PORT` no .env)

### Habilitação do gRPC
No arquivo `framework.yaml`:
```yaml
protocols:
  grpc: true
```

## 🧪 Testes Manuais

### Instalar grpcurl
```bash
# Go
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# macOS
brew install grpcurl

# Ubuntu/Debian
apt-get install grpcurl
```

### Exemplos de Comandos

#### Listar Serviços
```bash
grpcurl -plaintext localhost:50051 list
```

#### Health Check
```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

#### Criar Usuário
```bash
grpcurl -plaintext -d '{
  "name": "João Silva",
  "email": "joao@teste.com",
  "password": "123456"
}' localhost:50051 api.UserService/CreateUser
```

#### Criar Produto
```bash
grpcurl -plaintext -d '{
  "name": "Notebook Dell",
  "description": "Notebook Dell Inspiron 15",
  "price": 2500.99,
  "quantity": 10
}' localhost:50051 api.ProductService/CreateProduct
```

#### Listar Usuários
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 api.UserService/ListUsers
```

## 📁 Estrutura dos Arquivos

```
proto/                          # Definições Protocol Buffer
├── user.proto                 # Definições do serviço de usuário
├── product.proto              # Definições do serviço de produto
├── order.proto                # Definições do serviço de pedido
├── *.pb.go                    # Código Go gerado (não editar)
└── *_grpc.pb.go              # Código gRPC gerado (não editar)

internal/grpc/services/        # Implementações dos serviços gRPC
├── user_service.go            # Implementação UserService
├── product_service.go         # Implementação ProductService
└── order_service.go           # Implementação OrderService

pkg/framework/providers/grpc/  # Provider gRPC do framework
└── grpc.go                    # Configuração e inicialização do servidor gRPC
```

## 🔄 Regenerar Código Proto

Para regenerar o código Go a partir dos arquivos .proto:

```bash
make proto-gen
```

## 🏗️ Arquitetura

O gRPC está integrado ao framework modular:

1. **Providers**: O `GRPCProvider` gerencia o servidor gRPC
2. **Services**: Cada módulo tem seu próprio serviço gRPC
3. **Bootstrap**: Os serviços são registrados automaticamente no bootstrap
4. **Configuração**: Habilitado/desabilitado via `framework.yaml`

### Fluxo de Inicialização

1. Framework carrega configuração
2. Se `protocols.grpc` está habilitado:
   - Cria `GRPCProvider`
   - Inicializa servidor gRPC na porta configurada
   - Registra serviços dos módulos habilitados
   - Inicia servidor em goroutine separada

### Integração com Módulos Existentes

Os serviços gRPC reutilizam:
- **Service Layer**: Mesma lógica de negócio dos handlers HTTP
- **Repository Layer**: Mesma camada de persistência
- **Domain Logic**: Mesmas validações e regras
- **Events**: Mesma publicação de eventos

## 🔒 Considerações de Segurança

Por padrão, o gRPC roda sem TLS (plaintext) para desenvolvimento. Para produção:

1. Configure TLS no `GRPCProvider`
2. Implemente autenticação via interceptors
3. Configure rate limiting
4. Use network policies apropriadas

## 📊 Monitoramento

- Health checks disponíveis via `grpc.health.v1.Health/Check`
- Logs estruturados para todas as operações
- Métricas podem ser adicionadas via interceptors
- Reflection habilitada apenas em modo debug