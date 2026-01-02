package repositories

import (
	"context"
	"errors"
	"sync"

	"laba_7/domain"
)

var ErrOrderNotFound = errors.New("order not found")

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
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID()] = order
	return nil
}
