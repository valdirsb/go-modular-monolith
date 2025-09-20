# Configuração do Banco de Dados MySQL

Este documento descreve como configurar o banco de dados MySQL para o projeto.

## Opções de Configuração

### 🐳 Opção 1: MySQL com Docker (Recomendado)

- Docker e Docker Compose instalados
- Mais fácil e rápido de configurar

### 🖥️ Opção 2: MySQL Local

- MySQL Server 8.0 ou superior instalado
- Cliente MySQL (mysql-client ou MySQL Workbench)

## 🐳 Configuração com Docker (Recomendado)

### 1. Configuração Rápida

```bash
# Configuração inicial completa
make setup

# Ou passo a passo:
cp .env.example .env
make docker-up
```

### 2. Verificar se está funcionando

```bash
# Ver logs do MySQL
make docker-logs

# Conectar ao banco
make db-shell

# Testar a aplicação
make run
```

### 3. Comandos Úteis Docker

```bash
# Iniciar MySQL
make docker-up

# Parar MySQL  
make docker-down

# Ver logs
make docker-logs

# Acessar PhpMyAdmin
open http://localhost:8081
```

## 🖥️ Configuração Local (Alternativa)

### 1. Instalar MySQL (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install mysql-server
```

### 2. Configurar MySQL

```bash
# Iniciar o serviço MySQL
sudo systemctl start mysql

# Configurar a instalação (opcional)
sudo mysql_secure_installation
```

### 3. Criar o banco de dados

```bash
# Conectar ao MySQL como root
mysql -u root -p

# Executar o script de inicialização
source scripts/init_database.sql

# Ou executar manualmente:
CREATE DATABASE IF NOT EXISTS app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Criar usuário (opcional, se não quiser usar root)

```sql
-- Criar usuário específico para a aplicação
CREATE USER 'app_user'@'localhost' IDENTIFIED BY 'app_password';
GRANT ALL PRIVILEGES ON app_db.* TO 'app_user'@'localhost';
FLUSH PRIVILEGES;
```

## Configuração da Aplicação

As configurações de conexão estão no arquivo `internal/shared/database/database.go`:

```go
func GetDefaultConfig() *DatabaseConfig {
    return &DatabaseConfig{
        Host:     "localhost",
        Port:     "3306",
        Username: "root",
        Password: "123456",
        Database: "app_db",
    }
}
```

### Variáveis de Ambiente (Recomendado)

Para produção, use variáveis de ambiente:

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=123456
export DB_DATABASE=app_db
```

## Migrações

As migrações são executadas automaticamente na inicialização da aplicação através do GORM AutoMigrate.

As seguintes tabelas serão criadas:

- **users**: Dados dos usuários
  - id (VARCHAR(36) PRIMARY KEY)
  - username (VARCHAR(50) UNIQUE NOT NULL)
  - email (VARCHAR(100) UNIQUE NOT NULL)
  - password (VARCHAR(255) NOT NULL)
  - created_at (TIMESTAMP)
  - updated_at (TIMESTAMP)

## Testando a Conexão

Para testar se a conexão está funcionando:

```bash
# Execute a aplicação
go run cmd/server/main.go

# Você deve ver as seguintes mensagens:
# - "Successfully connected to MySQL database"
# - "Database migration completed successfully"
# - "Server starting on port 8080"
```

## Comandos Úteis

```bash
# Verificar status do MySQL
sudo systemctl status mysql

# Parar/Iniciar MySQL
sudo systemctl stop mysql
sudo systemctl start mysql

# Conectar ao banco
mysql -u root -p app_db

# Ver tabelas
USE app_db;
SHOW TABLES;
DESCRIBE users;

# Ver dados
SELECT * FROM users;
```

## Solução de Problemas

### Erro de Conexão

Se você receber erro de conexão:

1. Verifique se o MySQL está rodando:
   ```bash
   sudo systemctl status mysql
   ```

2. Verifique as credenciais no código

3. Teste a conexão manualmente:
   ```bash
   mysql -u root -p -h localhost -P 3306
   ```

### Erro de Permissão

Se houver erro de permissão:

```sql
GRANT ALL PRIVILEGES ON app_db.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
```

### Erro de Character Set

Se houver problemas com caracteres especiais:

```sql
ALTER DATABASE app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```