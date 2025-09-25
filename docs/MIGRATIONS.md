# Database Migrations

Este documento registra a evolução da estrutura do banco de dados.

## 🏗️ Database Schema Evolution

### Migration History

1. **Initial Setup (v1.0)**
   - Users table with UUID primary key
   - Basic user authentication fields

2. **Products Module (v1.1)** 
   - Products table with UUID primary key
   - Category support and advanced filtering
   - Stock management fields

3. **Orders Module (v1.2)**
   - Orders table with relationships
   - Order items for product quantities
   - Order status management
   - Stock integration with automatic updates

## 📊 Current Tables

### users
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

### products  
```sql
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INT DEFAULT 0,
    category_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_products_category (category_id),
    INDEX idx_products_price (price),
    INDEX idx_products_name (name)
);
```

### orders
```sql
CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    total DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_orders_user (user_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_created (created_at)
);
```

### order_items
```sql
CREATE TABLE order_items (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT,
    INDEX idx_order_items_order (order_id),
    INDEX idx_order_items_product (product_id)
);
```

## 🌱 Seeds Data

### Auto-Seeded Products (12 items)

| Category | Count | Price Range |
|----------|-------|-------------|
| electronics | 2 | R$ 8.999 - R$ 12.999 |
| computers | 2 | R$ 4.999 - R$ 12.999 |
| accessories | 2 | R$ 1.499 - R$ 2.299 |
| tablets | 2 | R$ 3.999 - R$ 8.999 |
| gaming | 2 | R$ 2.299 - R$ 4.999 |
| tv | 1 | R$ 15.999 |
| wearables | 1 | R$ 3.999 |

### Seed Details
```sql
-- Sample seeded products
INSERT INTO products VALUES
('prod-001', 'iPhone 15 Pro Max', 'Apple iPhone 15 Pro Max 256GB...', 8999.99, 15, 'electronics', NOW(), NOW()),
('prod-002', 'MacBook Air M2', 'MacBook Air 13" com chip M2...', 12999.99, 8, 'computers', NOW(), NOW()),
('prod-003', 'AirPods Pro', 'AirPods Pro (2ª geração)...', 2299.99, 25, 'accessories', NOW(), NOW()),
-- ... mais produtos
```

## 🔧 Database Operations

### Setup Commands
```bash
# Start database container
make docker-up

# Initialize application (auto-creates tables + seeds)
make run

# Access phpMyAdmin
open http://localhost:8081
```

### Manual Database Access
```bash
# Connect to MySQL container
docker exec -it mysql-container mysql -u root -p123456 app_db

# Show tables
SHOW TABLES;

# Check seeded products
SELECT id, name, category_id, price, stock FROM products ORDER BY category_id;
```

### Migration Rollback *(Se necessário)*
```sql
-- Remove all data and restart fresh
DROP DATABASE IF EXISTS app_db;
CREATE DATABASE app_db;
```

## 📈 Performance Considerations

### Current Indexes
- `users`: email, username (for login/uniqueness)
- `products`: category_id, price, name (for filtering)
- `orders`: user_id, status, created_at (for queries)
- `order_items`: order_id, product_id (for relationships)

### Query Optimization
- Product filters use compound indexes
- UUID primary keys for distributed-friendly IDs
- Timestamp indexes for date-range queries
- Foreign keys with appropriate cascade rules

## 🚀 Future Migrations

### Planned Features (v1.2)
1. **Order System Complete**
   - Order creation and management
   - Payment integration fields
   - Inventory reservation logic

2. **User Enhancements**
   - User profiles and preferences
   - Address management
   - Role-based permissions

3. **Product Enhancements**
   - Product images and galleries
   - Reviews and ratings system
   - Inventory tracking history

### Database Growth Planning
- Partitioning strategy for large tables
- Archive tables for historical data
- Read replicas for query performance
- Backup and recovery procedures