package domain

type OrderRepository interface {
	Save(order *Order) error
	FindByID(id string) (*Order, error)
	FindAll() ([]*Order, error)
	Delete(id string) error
	Update(order *Order) error
}
