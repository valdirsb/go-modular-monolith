# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
e este projeto segue [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-09-23

### ✨ Adicionado
- **Order Module**: Sistema completo de pedidos
  - Criação de pedidos com múltiplos itens
  - Gestão de status (pending, confirmed, shipped, delivered, cancelled)
  - Integração automática com controle de estoque
  - Cancelamento com reversão de estoque
  - Eventos de domínio (OrderCreated, OrderStatusUpdated, OrderCancelled)
  
- **Otimizações de Performance**:
  - Cache de produtos no OrderService para evitar consultas repetidas
  - Agregação de quantidades por produto
  - Operações transacionais para consistência
  
- **Documentação Completa**:
  - README específico para cada módulo
  - Documentação detalhada de arquitetura
  - API endpoints com exemplos
  - Guias de desenvolvimento e deployment

- **Testes Extensivos**:
  - Testes de integração para todos os módulos
  - Cenários de erro e validação
  - Script automatizado de testes da API

### 🔧 Melhorado
- **API Documentation**: Endpoints do Order module adicionados
- **Database Migrations**: Tabelas orders e order_items
- **Test Script**: Testes completos do fluxo de pedidos
- **Makefile**: Novos comandos para docs e estatísticas

### 🐛 Corrigido
- Correção na lógica de atualização de estoque (total vs delta)
- Validações de negócio mais robustas
- Tratamento de erros aprimorado

## [1.1.0] - 2025-09-20

### ✨ Adicionado
- **Product Module**: Sistema completo de catálogo
  - CRUD de produtos com validações
  - Sistema de categorias
  - Controle de estoque
  - Filtros avançados (categoria, preço, nome)
  - Paginação para grandes volumes
  - 12 produtos seedados em 7 categorias
  
- **Events System**: 
  - EventBus para comunicação entre módulos
  - ProductCreatedEvent e ProductStockUpdatedEvent
  - Infraestrutura para eventos assíncronos

- **API Enhancements**:
  - Query parameters para filtros
  - Endpoint específico para atualização de estoque
  - Respostas padronizadas com códigos HTTP

### 🔧 Melhorado
- **Database**: Índices otimizados para performance
- **Architecture**: Separação clara entre domínio e infraestrutura
- **Documentation**: API docs com exemplos práticos

### 📊 Dados Seedados
- iPhone 15 Pro Max, Samsung Galaxy S24 Ultra (electronics)
- MacBook Air M2, Dell XPS 13 (computers)  
- AirPods Pro, Sony WH-1000XM5 (accessories)
- iPad Air, Microsoft Surface Pro 9 (tablets)
- Nintendo Switch OLED, PlayStation 5 (gaming)
- LG OLED C3 55" (tv)
- Apple Watch Series 9 (wearables)

## [1.0.0] - 2025-09-15

### ✨ Primeira Versão
- **User Module**: Sistema completo de usuários
  - CRUD de usuários com validações
  - Autenticação com hash Argon2
  - Validação de credenciais
  - UserCreatedEvent

- **Infrastructure Foundation**:
  - Clean Architecture + Hexagonal Architecture + DDD
  - Dependency Injection Container
  - MySQL com GORM
  - Docker Compose para desenvolvimento
  - Gin HTTP framework

- **Development Tools**:
  - Makefile com comandos essenciais
  - Docker setup automatizado
  - Health check endpoint
  - Estrutura modular escalável

### 🏗️ Arquitetura Base
- **Modular Monolith**: Módulos independentes e desacoplados
- **Clean Architecture**: Dependências sempre para o domínio
- **DDD**: Modelagem rica com agregados e entidades
- **Ports & Adapters**: Interfaces bem definidas
- **Event-Driven**: Comunicação assíncrona entre módulos

---

## 🚀 Próximas Versões Planejadas

### [1.3.0] - Q4 2025
- [ ] **Payment Module**: Integração com gateways de pagamento
- [ ] **Notification Module**: Email e SMS
- [ ] **Authentication**: JWT tokens e sessões
- [ ] **Authorization**: Roles e permissões

### [1.4.0] - Q1 2026  
- [ ] **Inventory Module**: Gestão avançada de estoque
- [ ] **Shipping Module**: Cálculo de frete e rastreamento
- [ ] **Analytics Module**: Relatórios e métricas
- [ ] **API Versioning**: Suporte a múltiplas versões

### [2.0.0] - Q2 2026
- [ ] **Microservices Migration**: Preparação para split
- [ ] **GraphQL API**: Alternativa ao REST
- [ ] **Event Sourcing**: Histórico completo de eventos
- [ ] **CQRS**: Separação de leitura e escrita

---

## 📋 Tipos de Mudanças

- **✨ Adicionado**: Para novas funcionalidades
- **🔧 Melhorado**: Para mudanças em funcionalidades existentes  
- **🐛 Corrigido**: Para correções de bugs
- **📚 Documentação**: Para mudanças apenas na documentação
- **🔒 Segurança**: Para correções de vulnerabilidades
- **💥 Removido**: Para funcionalidades removidas
- **⚠️ Deprecated**: Para funcionalidades que serão removidas

## 🔗 Links

- [Repositório GitHub](https://github.com/valdirsb/go-modular-monolith)
- [Documentação API](./docs/API.md)
- [Guia de Arquitetura](./ARCHITECTURE.md)
- [Guia de Desenvolvimento](./docs/DEVELOPMENT.md)