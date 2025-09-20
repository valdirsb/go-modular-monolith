package domain

type Order struct {
    ID          string
    UserID      string
    ProductID   string
    Quantity    int
    Status      string
}

func NewOrder(id, userID, productID string, quantity int, status string) *Order {
    return &Order{
        ID:        id,
        UserID:    userID,
        ProductID: productID,
        Quantity:  quantity,
        Status:    status,
    }
}

func (o *Order) UpdateStatus(newStatus string) {
    o.Status = newStatus
}