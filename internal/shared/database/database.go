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

// ProductModel representa a estrutura da tabela products no banco
type ProductModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:500"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Stock       int       `gorm:"default:0;not null"`
	CategoryID  string    `gorm:"size:36;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela
func (ProductModel) TableName() string {
	return "products"
}

// ToContract converte ProductModel para contracts.Product
func (p *ProductModel) ToContract() *contracts.Product {
	return &contracts.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// FromContract converte contracts.Product para ProductModel
func (p *ProductModel) FromContract(product *contracts.Product) {
	p.ID = product.ID
	p.Name = product.Name
	p.Description = product.Description
	p.Price = product.Price
	p.Stock = product.Stock
	p.CategoryID = product.CategoryID
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt
}

// OrderModel representa a estrutura da tabela orders no banco
type OrderModel struct {
	ID        string           `gorm:"primaryKey;size:36"`
	UserID    string           `gorm:"size:36;not null;index"`
	Status    string           `gorm:"size:20;not null"`
	Total     float64          `gorm:"type:decimal(10,2);not null"`
	Items     []OrderItemModel `gorm:"foreignKey:OrderID"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela
func (OrderModel) TableName() string {
	return "orders"
}

// OrderItemModel representa a estrutura da tabela order_items no banco
type OrderItemModel struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   string  `gorm:"size:36;not null;index"`
	ProductID string  `gorm:"size:36;not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
}

// TableName especifica o nome da tabela
func (OrderItemModel) TableName() string {
	return "order_items"
}

// ToContract converte OrderModel para contracts.Order
func (o *OrderModel) ToContract() *contracts.Order {
	items := make([]contracts.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = contracts.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &contracts.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		Items:     items,
		Status:    contracts.OrderStatus(o.Status),
		Total:     o.Total,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

// FromContract converte contracts.Order para OrderModel
func (o *OrderModel) FromContract(order *contracts.Order) {
	o.ID = order.ID
	o.UserID = order.UserID
	o.Status = string(order.Status)
	o.Total = order.Total
	o.CreatedAt = order.CreatedAt
	o.UpdatedAt = order.UpdatedAt

	// Converter items
	o.Items = make([]OrderItemModel, len(order.Items))
	for i, item := range order.Items {
		o.Items[i] = OrderItemModel{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
}

// AutoMigrate executa as migrações necessárias
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&UserModel{},
		&ProductModel{},
		&OrderModel{},
		&OrderItemModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	log.Println("Database migration completed successfully")

	// Executar seeds
	if err := SeedDatabase(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
		// Não falhar se o seed der erro, apenas avisar
	}

	return nil
}

// SeedDatabase popula o banco com dados iniciais
func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Verificar se já existem produtos para evitar duplicação
	var count int64
	if err := db.Model(&ProductModel{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count products: %w", err)
	}

	if count > 0 {
		log.Printf("Database already has %d products, skipping seed", count)
		return nil
	}

	// Seeds de produtos
	products := []ProductModel{
		{
			ID:          "prod-001",
			Name:        "iPhone 15 Pro Max",
			Description: "Apple iPhone 15 Pro Max 256GB - Titânio Natural com câmera profissional de 48MP",
			Price:       8999.99,
			Stock:       15,
			CategoryID:  "electronics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-002",
			Name:        "MacBook Air M2",
			Description: "MacBook Air 13\" com chip M2, 8GB RAM, 256GB SSD - Cor Meia-noite",
			Price:       12999.99,
			Stock:       8,
			CategoryID:  "computers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-003",
			Name:        "Samsung Galaxy S24 Ultra",
			Description: "Samsung Galaxy S24 Ultra 512GB - Preto com S Pen incluída e câmera de 200MP",
			Price:       7499.99,
			Stock:       12,
			CategoryID:  "electronics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-004",
			Name:        "Dell XPS 13",
			Description: "Notebook Dell XPS 13 Intel Core i7, 16GB RAM, 512GB SSD, Tela InfinityEdge",
			Price:       9999.99,
			Stock:       6,
			CategoryID:  "computers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-005",
			Name:        "AirPods Pro (3ª geração)",
			Description: "Apple AirPods Pro com cancelamento ativo de ruído e case de carregamento MagSafe",
			Price:       2499.99,
			Stock:       25,
			CategoryID:  "accessories",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-006",
			Name:        "Sony WH-1000XM5",
			Description: "Fone de ouvido Sony WH-1000XM5 com cancelamento de ruído premium e 30h de bateria",
			Price:       1899.99,
			Stock:       18,
			CategoryID:  "accessories",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-007",
			Name:        "iPad Air (5ª geração)",
			Description: "iPad Air 10.9\" com chip M1, 256GB, Wi-Fi + Cellular - Azul-céu",
			Price:       6499.99,
			Stock:       10,
			CategoryID:  "tablets",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-008",
			Name:        "Nintendo Switch OLED",
			Description: "Console Nintendo Switch modelo OLED com tela de 7\" e 64GB de armazenamento",
			Price:       2799.99,
			Stock:       20,
			CategoryID:  "gaming",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-009",
			Name:        "PlayStation 5",
			Description: "Console Sony PlayStation 5 com SSD ultrarrápido e controle DualSense",
			Price:       4999.99,
			Stock:       5,
			CategoryID:  "gaming",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-010",
			Name:        "Microsoft Surface Pro 9",
			Description: "Surface Pro 9 Intel i7, 16GB RAM, 512GB SSD com teclado Type Cover incluso",
			Price:       11999.99,
			Stock:       7,
			CategoryID:  "tablets",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-011",
			Name:        "LG OLED C3 55\"",
			Description: "Smart TV LG OLED C3 55\" 4K com webOS, Dolby Vision IQ e Gaming Hub",
			Price:       6999.99,
			Stock:       4,
			CategoryID:  "tv",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-012",
			Name:        "Apple Watch Series 9",
			Description: "Apple Watch Series 9 GPS 45mm caixa de alumínio com pulseira esportiva",
			Price:       3999.99,
			Stock:       14,
			CategoryID:  "wearables",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Inserir produtos em lote
	if err := db.CreateInBatches(products, 100).Error; err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	log.Printf("Successfully seeded %d products", len(products))
	return nil
}
