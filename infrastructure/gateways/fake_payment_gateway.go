package gateways

import (
	"context"

	"laba_7/domain"
)

type FakePaymentGateway struct{}

func NewFakePaymentGateway() *FakePaymentGateway {
	return &FakePaymentGateway{}
}

func (g *FakePaymentGateway) Charge(ctx context.Context, orderID domain.OrderID, amount domain.Money) error {
	if amount.Amount() <= 0 {
		return domain.ErrEmptyOrder
	}
	return nil
}
