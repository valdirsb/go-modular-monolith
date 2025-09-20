package domain

type UserRepository interface {
    Create(user *User) error
    GetByID(id string) (*User, error)
    Update(user *User) error
    Delete(id string) error
    GetAll() ([]*User, error)
}