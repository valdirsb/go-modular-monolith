package domain

import (
	"errors"
	"regexp"
	"time"
	"unicode/utf8"

	"go-modular-monolith/pkg/contracts"
)

// User representa a entidade de domínio do usuário
type User struct {
	contracts.User
}

// UserAggregate contém as regras de negócio do usuário
type UserAggregate struct {
	user *User
}

// NewUser cria um novo usuário com validações de domínio
func NewUser(id, username, email string) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	return &User{
		User: contracts.User{
			ID:        id,
			Username:  username,
			Email:     email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// NewUserAggregate cria um novo aggregate de usuário
func NewUserAggregate(user *User) *UserAggregate {
	return &UserAggregate{user: user}
}

// UpdateEmail atualiza o email do usuário com validação
func (ua *UserAggregate) UpdateEmail(newEmail string) error {
	if err := validateEmail(newEmail); err != nil {
		return err
	}

	ua.user.Email = newEmail
	ua.user.UpdatedAt = time.Now()
	return nil
}

// UpdateUsername atualiza o username do usuário com validação
func (ua *UserAggregate) UpdateUsername(newUsername string) error {
	if err := validateUsername(newUsername); err != nil {
		return err
	}

	ua.user.Username = newUsername
	ua.user.UpdatedAt = time.Now()
	return nil
}

// SetPassword define uma nova senha (hash)
func (ua *UserAggregate) SetPassword(hashedPassword string) {
	ua.user.Password = hashedPassword
	ua.user.UpdatedAt = time.Now()
}

// GetUser retorna o usuário do aggregate
func (ua *UserAggregate) GetUser() *User {
	return ua.user
}

// IsValid verifica se o usuário é válido
func (ua *UserAggregate) IsValid() error {
	if err := validateUsername(ua.user.Username); err != nil {
		return err
	}

	if err := validateEmail(ua.user.Email); err != nil {
		return err
	}

	if ua.user.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	return nil
}

// Domain validation functions

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	if utf8.RuneCountInString(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if utf8.RuneCountInString(username) > 50 {
		return errors.New("username must be at most 50 characters long")
	}

	// Apenas letras, números, underscore e hífen
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers, underscore and hyphens")
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	// Regex básico para validação de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}
