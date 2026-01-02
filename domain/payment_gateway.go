package domain

import "context"

type PaymentGateway interface {
	Charge(ctx context.Context, orderID OrderID, amount Money) error
}
