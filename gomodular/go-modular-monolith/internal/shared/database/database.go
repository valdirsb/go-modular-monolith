package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-modular-monolith/pkg/contracts"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig contém as configurações de conexão do banco
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// GetDefaultConfig retorna a configuração do banco usando variáveis de ambiente ou padrões
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		Username: getEnvOrDefault("DB_USERNAME", "root"),
		Password: getEnvOrDefault("DB_PASSWORD", "123456"),
		Database: getEnvOrDefault("DB_DATABASE", "app_db"),
	}
}

// getEnvOrDefault obtém valor de variável de ambiente ou retorna padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect estabelece conexão com o banco MySQL
func Connect(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configurações do pool de conexões
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to MySQL database")
	return db, nil
}

// UserModel representa a estrutura da tabela users no banco
type UserModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	Email     string    `gorm:"uniqueIndex;size:100;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela
func (UserModel) TableName() string {
	return "users"
}

// ToContract converte UserModel para contracts.User
func (u *UserModel) ToContract() *contracts.User {
	return &contracts.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromContract converte contracts.User para UserModel
func (u *UserModel) FromContract(user *contracts.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.Email = user.Email
	u.Password = user.Password
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

// AutoMigrate executa as migrações necessárias
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&UserModel{},
		// Adicionar outros modelos aqui conforme necessário
	)
	if err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
