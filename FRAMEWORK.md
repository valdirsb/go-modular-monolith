# Go Modular Framework

Um **framework modular** em Go que permite criar sistemas escaláveis e configuráveis através de um arquivo de configuração. Ative ou desative funcionalidades conforme necessário para cada projeto.

## 🎯 Características do Framework

### 🏗️ Arquitetura Flexível
- **Modular**: Cada funcionalidade é um módulo independente que pode ser ativado/desativado
- **Configurável**: Todas as funcionalidades são controladas via `framework.yaml`
- **Extensível**: Adicione novos providers e integrações facilmente
- **Clean Architecture**: Separação clara entre domínio, aplicação e infraestrutura

### 📦 Funcionalidades Incluídas

#### Core Features
- ✅ **HTTP Server**: Servidor web com Gin
- ✅ **Dependency Injection**: Container DI personalizado
- ✅ **Event System**: Comunicação entre módulos via eventos
- ✅ **Logging**: Sistema de logs estruturado
- ✅ **Database Migrations**: Migrações automáticas

#### Database Providers
- ✅ **MySQL**: Suporte completo via GORM
- 🚧 **PostgreSQL**: Em desenvolvimento
- 🚧 **MongoDB**: Em desenvolvimento
- 🚧 **SQLite**: Em desenvolvimento

#### Cache Providers
- 🚧 **Redis**: Em desenvolvimento
- 🚧 **Memcached**: Em desenvolvimento
- 🚧 **In-Memory**: Em desenvolvimento

#### Communication Protocols
- ✅ **HTTP**: REST APIs
- 🚧 **gRPC**: Comunicação high-performance
- 🚧 **WebSockets**: Comunicação real-time
- 🚧 **GraphQL**: APIs flexíveis

#### External Integrations
- ✅ **Payment Gateways**: Stripe, PIX, PayPal, Mercado Pago
- ✅ **Maps & Location**: Google Maps, Mapbox
- 🚧 **Cloud Storage**: AWS S3, Google Cloud, Azure
- 🚧 **Authentication**: Firebase Auth, Auth0, OAuth2
- 🚧 **Messaging**: RabbitMQ, Kafka, SQS
- 🚧 **Email**: SendGrid, SES, SMTP
- 🚧 **Push Notifications**: Firebase Messaging, APNS
- 🚧 **Monitoring**: Prometheus, Jaeger, DataDog

#### Security Features
- ✅ **CORS**: Cross-Origin Resource Sharing
- ✅ **JWT**: JSON Web Tokens
- 🚧 **Rate Limiting**: Controle de taxa de requisições
- 🚧 **Encryption**: Criptografia de dados
- 🚧 **CSRF Protection**: Proteção contra CSRF

## 🚀 Início Rápido

### 1. Configuração do Framework
```yaml
# framework.yaml
framework:
  name: "meu-projeto"
  version: "1.0.0"

# Ative apenas o que precisar
database:
  mysql: true
  redis: false

integrations:
  stripe: true
  google_maps: false
  
modules:
  user: true
  product: true
  order: false
```

### 2. Comandos Disponíveis
```bash
# Configuração inicial
make setup-framework

# Executar aplicação
make run-framework

# Informações do framework
make framework-info

# Health check
curl http://localhost:8080/health

# Framework info
curl http://localhost:8080/framework/info
```

### 3. Criar Novo Projeto
```bash
# Criar projeto básico
./scripts/create_project.sh meu-projeto basic

# Criar API completa
./scripts/create_project.sh minha-api api

# Criar sistema de pagamentos
./scripts/create_project.sh pagamentos payment

# Criar com todas funcionalidades
./scripts/create_project.sh projeto-full full
```

## 📋 Templates Disponíveis

### Basic Template
Projeto mínimo com HTTP + MySQL
- Módulos: User, Product
- Database: MySQL
- Protocolos: HTTP
- Integrações: JWT

### API Template
API REST completa
- Módulos: User, Product, Order, Notification
- Database: MySQL + Redis
- Protocolos: HTTP
- Integrações: JWT, OAuth2
- Security: CORS, Rate Limiting, CSRF
- Development: Swagger

### Payment Template
Sistema de pagamentos
- Módulos: User, Billing, Order
- Database: MySQL + Redis
- Integrações: Stripe, PayPal, PIX, Mercado Pago
- Security: CORS, Rate Limiting, Encryption

### Maps Template
Sistema com funcionalidades de mapas
- Módulos: User, Location
- Database: MySQL
- Integrações: Google Maps

### Full Template
Todas as funcionalidades ativadas
- Todos os módulos
- Todos os bancos de dados
- Todas as integrações
- Todas as funcionalidades de segurança

## 🏗️ Estrutura do Projeto

```
projeto/
├── cmd/server/
│   ├── main.go              # Aplicação tradicional
│   └── main_framework.go    # Aplicação com framework
├── internal/
│   ├── bootstrap/
│   │   ├── bootstrap.go           # Bootstrap tradicional
│   │   └── framework_bootstrap.go # Bootstrap com framework
│   └── modules/             # Módulos de negócio
├── pkg/
│   ├── framework/           # Core do framework
│   │   ├── config.go       # Configuração
│   │   ├── framework.go    # Framework principal
│   │   └── providers/      # Providers de tecnologias
├── framework.yaml           # Configuração do framework
└── Makefile                # Comandos de desenvolvimento
```

## 🔧 Extensibilidade

### Adicionando Novo Provider
```go
// pkg/framework/providers/myprovider/provider.go
type MyProvider struct {
    config *framework.FrameworkConfig
}

func (p *MyProvider) Name() string {
    return "my_provider"
}

func (p *MyProvider) IsEnabled() bool {
    return framework.IsEnabled("integrations", "my_service")
}

func (p *MyProvider) Initialize(ctx context.Context, deps providers.Dependencies) error {
    // Inicialização do provider
    return nil
}
```

### Adicionando Novo Módulo
```bash
# Criar estrutura do módulo
mkdir -p internal/modules/meu_modulo/{domain,service,handler,repository}

# Adicionar no framework.yaml
modules:
  meu_modulo: true

# Registrar no bootstrap
if framework.IsEnabled("modules", "meu_modulo") {
    registerMeuModulo(c)
}
```

## 📊 Monitoramento

### Health Check
```bash
curl http://localhost:8080/health
```
```json
{
  "status": "healthy",
  "services": {
    "stripe": null,
    "mysql": null
  },
  "features": {
    "core": ["http_server", "event_system"],
    "integrations": ["stripe"],
    "modules": ["user", "product"]
  }
}
```

### Framework Info
```bash
curl http://localhost:8080/framework/info
```
```json
{
  "framework": {
    "name": "meu-projeto",
    "version": "1.0.0"
  },
  "features": {
    "core": ["http_server"],
    "database": ["mysql"],
    "modules": ["user", "product"]
  }
}
```

## 🛠️ Desenvolvimento

### Estrutura de Desenvolvimento
1. **Configure** o que precisar no `framework.yaml`
2. **Execute** `make run-framework`
3. **Desenvolva** seus módulos de negócio
4. **Teste** via endpoints de health e framework info
5. **Deploy** apenas com as funcionalidades ativas

### Boas Práticas
- Mantenha módulos independentes
- Use o sistema de eventos para comunicação entre módulos
- Configure apenas o que vai usar para reduzir complexidade
- Monitore via health checks
- Use templates para projetos novos

## 📚 Exemplos de Uso

### Sistema de E-commerce
```yaml
modules:
  user: true
  product: true
  order: true
  inventory: true
  
integrations:
  stripe: true
  sendgrid: true
  
security:
  cors: true
  rate_limiting: true
```

### Sistema de Delivery
```yaml
modules:
  user: true
  order: true
  location: true
  
integrations:
  google_maps: true
  mercado_pago: true
  firebase_messaging: true
```

### Sistema Corporativo
```yaml
modules:
  user: true
  analytics: true
  audit: true
  
integrations:
  oauth2: true
  prometheus: true
  
security:
  cors: true
  rate_limiting: true
  encryption: true
```

---

## 🤝 Contribuição

Contribuições são bem-vindas! Adicione novos providers, módulos ou melhorias.

## 📄 Licença

MIT License