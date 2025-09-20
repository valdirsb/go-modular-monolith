package repository

import (
	"context"
	"fmt"

	"go-modular-monolith/internal/modules/user/ports"
	"go-modular-monolith/internal/shared/database"
	"go-modular-monolith/pkg/contracts"

	"gorm.io/gorm"
)

// mysqlUserRepository implementa a interface UserRepository usando MySQL/GORM
type mysqlUserRepository struct {
	db *gorm.DB
}

// NewMySQLUserRepository cria uma nova instância do repositório MySQL
func NewMySQLUserRepository(db *gorm.DB) ports.UserRepository {
	return &mysqlUserRepository{
		db: db,
	}
}

// Create cria um novo usuário no banco de dados
func (r *mysqlUserRepository) Create(ctx context.Context, user *contracts.User) error {
	userModel := &database.UserModel{}
	userModel.FromContract(user)

	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID busca um usuário pelo ID
func (r *mysqlUserRepository) GetByID(ctx context.Context, id string) (*contracts.User, error) {
	var userModel database.UserModel

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return userModel.ToContract(), nil
}

// GetByEmail busca um usuário pelo email
func (r *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (*contracts.User, error) {
	var userModel database.UserModel

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return userModel.ToContract(), nil
}

// Update atualiza um usuário existente
func (r *mysqlUserRepository) Update(ctx context.Context, user *contracts.User) error {
	userModel := &database.UserModel{}
	userModel.FromContract(user)

	result := r.db.WithContext(ctx).Where("id = ?", user.ID).Updates(userModel)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete remove um usuário do banco de dados
func (r *mysqlUserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&database.UserModel{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
