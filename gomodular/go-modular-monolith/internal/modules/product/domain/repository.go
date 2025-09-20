package domain

type ProductRepository interface {
	FindByID(id string) (*Product, error)
	Save(product *Product) error
	Delete(id string) error
}
