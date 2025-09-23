# User Module 👤

O módulo User implementa o gerenciamento completo de usuários seguindo os princípios de Clean Architecture e DDD.

## 🎯 Responsabilidades

- **Gestão de Usuários**: CRUD completo (Create, Read, Update, Delete)
- **Autenticação Segura**: Hash de senhas com Argon2
- **Validação de Credenciais**: Login e verificação de senhas
- **Eventos de Domínio**: Notificação quando usuários são criados

## 🏗️ Estrutura do Módulo

```
user/
├── domain/
│   ├── user.go           # Entidade User e agregado UserAggregate
│   └── repository.go     # Interface UserRepository (Port)
├── service/
│   └── user_service.go   # Casos de uso e lógica de aplicação
├── repository/
│   ├── user_repository.go      # Interface específica
│   └── mysql_user_repository.go # Implementação MySQL
├── handler/
│   └── user_handler.go   # HTTP handlers
├── adapters/
│   └── password_hasher.go # Hash de senhas com Argon2
└── ports/
    └── ports.go          # Interfaces específicas do módulo
```

## 📡 API Endpoints

### Base URL: `/api/v1/users/`

#### Criar Usuário
```http
POST /api/v1/users/
Content-Type: application/json

{
  "username": "joao123",
  "email": "joao@example.com",
  "password": "senha123456"
}
```

**Response (201):**
```json
{
  "id": "uuid-generated",
  "username": "joao123",
  "email": "joao@example.com",
  "created_at": "2025-09-23T20:33:20Z",
  "updated_at": "2025-09-23T20:33:20Z"
}
```

#### Buscar Usuário
```http
GET /api/v1/users/{id}
```

#### Atualizar Usuário
```http
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "username": "joao_updated"
}
```

#### Deletar Usuário
```http
DELETE /api/v1/users/{id}
```

#### Validar Credenciais
```http
POST /api/v1/users/validate
Content-Type: application/json

{
  "email": "joao@example.com",
  "password": "senha123456"
}
```

## 🔒 Segurança

### Hash de Senhas
- **Algoritmo**: Argon2id (padrão da indústria)
- **Configuração**: Otimizada para segurança e performance
- **Salt**: Único por senha, gerado automaticamente

### Validações
- **Email**: Formato válido e único
- **Username**: Mínimo 3 caracteres, único
- **Password**: Mínimo 6 caracteres

## 🔄 Eventos Publicados

### UserCreatedEvent
Disparado quando um novo usuário é criado com sucesso.

```json
{
  "type": "UserCreatedEventType",
  "payload": {
    "user_id": "uuid",
    "email": "user@example.com"
  },
  "timestamp": "2025-09-23T20:33:20Z"
}
```

## 🏛️ Arquitetura

### Camada de Domínio
- **User**: Entidade com validações básicas
- **UserAggregate**: Agregado com regras de negócio complexas
- **UserRepository**: Interface para persistência

### Camada de Aplicação
- **UserService**: Orquestra casos de uso
  - Criação de usuários com hash de senha
  - Validação de credenciais
  - Publicação de eventos

### Camada de Infraestrutura
- **MySQLUserRepository**: Persistência em MySQL
- **Argon2PasswordHasher**: Hash seguro de senhas
- **UserHandler**: Endpoints HTTP REST

## 🗄️ Modelo de Dados

### Tabela: `users`
```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_username (username)
);
```

## 🧪 Testes

O módulo User é testado através do script `test_api.sh`:

1. **Criação de usuário**: com dados únicos gerados por timestamp
2. **Busca por ID**: verificação dos dados criados
3. **Atualização**: modificação de username
4. **Validação de credenciais**: login com senha correta
5. **Validação com senha incorreta**: deve falhar
6. **Exclusão**: remoção do usuário
7. **Busca após exclusão**: deve retornar 404

### Executar Testes
```bash
# Todos os testes
./scripts/test_api.sh

# Apenas usuários (manual)
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'
```

## ⚡ Performance

### Otimizações Implementadas
- **Índices**: email e username para busca rápida
- **UUID**: IDs únicos distribuídos
- **Hash Eficiente**: Argon2 com parâmetros otimizados

### Métricas Esperadas
- **Criação de usuário**: < 100ms
- **Busca por ID**: < 10ms
- **Validação de credenciais**: < 50ms (devido ao hash)

## 🚨 Limitações e Melhorias Futuras

### Limitações Atuais
- Não há sistema de roles/permissões
- Não há recuperação de senha
- Não há bloqueio de conta por tentativas

### Roadmap
- [ ] Sistema de roles (admin, user, etc.)
- [ ] Reset de senha via email
- [ ] Bloqueio temporário por tentativas incorretas
- [ ] JWT tokens para sessões
- [ ] 2FA (autenticação de dois fatores)
- [ ] Auditoria de login (logs de acesso)