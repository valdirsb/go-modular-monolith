# Módulos de Domínio

Este diretório contém todos os módulos de domínio da aplicação, organizados seguindo os princípios de Clean Architecture e DDD.

## 📁 Estrutura Padrão de Módulos

Cada módulo segue a mesma estrutura organizacional:

```
{module}/
├── domain/              # Entidades de domínio e regras de negócio
│   ├── {entity}.go      # Entidades e agregados do domínio
│   └── repository.go    # Interface do repositório (Port)
├── service/             # Casos de uso e lógica de aplicação
│   └── {module}_service.go
├── repository/          # Implementação de persistência (Adapter)
│   └── mysql_{module}_repository.go
├── handler/             # Controllers HTTP (Adapter)
│   └── {module}_handler.go
├── adapters/           # Adaptadores para serviços externos (opcional)
│   └── {adapter}.go
└── ports/              # Interfaces específicas do módulo (opcional)
    └── ports.go
```

## 🏗️ Módulos Implementados

### [User Module](./user/) 👤
Gerenciamento completo de usuários com autenticação segura.

**Funcionalidades:**
- CRUD de usuários
- Hash de senha com Argon2
- Validação de credenciais
- Eventos de usuário criado

### [Product Module](./product/) 📦
Sistema de catálogo de produtos com controle de estoque.

**Funcionalidades:**
- CRUD de produtos
- Controle de estoque
- Filtros avançados (categoria, preço, nome)
- Paginação
- Eventos de produto e estoque

### [Order Module](./order/) 🛒
Sistema completo de pedidos com gestão de estoque integrada.

**Funcionalidades:**
- Criação de pedidos com múltiplos itens
- Validação automática de estoque
- Gestão de status do pedido
- Cancelamento com reversão de estoque
- Cache de produtos para performance
- Eventos de pedido (criado, atualizado, cancelado)

## 🔧 Princípios Arquiteturais

### Separação de Responsabilidades
- **Domain**: Regras de negócio puras, sem dependências externas
- **Service**: Orquestração de casos de uso, coordena domain + infrastructure
- **Repository**: Persistência de dados, implementa interfaces do domain
- **Handler**: Entrada HTTP, converte requests em chamadas de service
- **Adapters**: Integrações com serviços externos (email, pagamento, etc.)

### Inversão de Dependência
- Services dependem de interfaces (repositories, adapters)
- Domain não conhece infraestrutura
- Handlers delegam para services
- DI container resolve dependências

### Comunicação Entre Módulos
- **Events**: Comunicação assíncrona via EventBus
- **Interfaces**: Chamadas diretas entre services quando necessário
- **Evitar**: Acoplamento direto entre repositories

## 📋 Checklist para Novo Módulo

Ao criar um novo módulo, certifique-se de:

- [ ] Definir entidades no `domain/`
- [ ] Criar interfaces no `contracts/interfaces.go`
- [ ] Implementar service com lógica de negócio
- [ ] Implementar repository MySQL
- [ ] Criar handler HTTP
- [ ] Registrar no bootstrap DI
- [ ] Adicionar rotas no main.go
- [ ] Incluir na migração automática
- [ ] Criar testes no script de API
- [ ] Documentar endpoints na API.md
- [ ] Publicar eventos quando apropriado

## 🚀 Como Executar

```bash
# Setup completo
make setup

# Executar aplicação
make run

# Testar todos os módulos
./scripts/test_api.sh
```

## 🧪 Testes

Cada módulo possui testes integrados no script `scripts/test_api.sh`:

- **User**: Criação, busca, atualização, validação, exclusão
- **Product**: CRUD, filtros, estoque, categorias
- **Order**: Criação, status, cancelamento, validações de estoque

Execute com:
```bash
./scripts/test_api.sh
```