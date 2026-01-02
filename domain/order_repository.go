package domain

import "context"

type OrderRepository interface {
	GetByID(ctx context.Context, id OrderID) (*Order, error)
	Save(ctx context.Context, order *Order) error
}
