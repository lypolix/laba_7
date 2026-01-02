package repositories

import (
	"context"
	"sync"

	"laba_7/domain"
)

type InMemoryOrderRepository struct {
	orders map[domain.OrderID]*domain.Order
	mu     sync.RWMutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[domain.OrderID]*domain.Order),
	}
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	r.mu.RLock()
	order, exists := r.orders[id]
	r.mu.RUnlock()

	if !exists {
		lines := []domain.OrderLine{
			{
				ProductID: "prod1",
				Quantity:  2,
				Price:     domain.NewMoney(5000, "RUB"),
			},
			{
				ProductID: "prod2",
				Quantity:  1,
				Price:     domain.NewMoney(10000, "RUB"),
			},
		}
		order, err := domain.NewOrder(id, lines)
		if err != nil {
			return nil, err
		}
		return order, nil
	}

	lines := make([]domain.OrderLine, len(order.Lines()))
	copy(lines, order.Lines())
	newOrder, err := domain.NewOrder(id, lines)
	if err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if err := order.ValidateTotal(); err != nil {
		return err
	}
	
	r.orders[order.ID()] = order
	return nil
}
