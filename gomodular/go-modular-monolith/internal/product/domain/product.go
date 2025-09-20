package domain

type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

func NewProduct(id, name, description string, price float64) *Product {
    return &Product{
        ID:          id,
        Name:        name,
        Description: description,
        Price:       price,
    }
}

func (p *Product) UpdateName(name string) {
    p.Name = name
}

func (p *Product) UpdateDescription(description string) {
    p.Description = description
}

func (p *Product) UpdatePrice(price float64) {
    p.Price = price
}