package application

import (
	"context"
	"laba_7/domain"
)

type PayOrderResult struct {
	OrderID domain.OrderID
	Success bool
	Error   error
}

type PayOrderUseCase struct {
	orderRepo domain.OrderRepository
	paymentGW domain.PaymentGateway
}

func NewPayOrderUseCase(repo domain.OrderRepository, gw domain.PaymentGateway) *PayOrderUseCase {
	return &PayOrderUseCase{
		orderRepo: repo,
		paymentGW: gw,
	}
}

func (uc *PayOrderUseCase) Execute(ctx context.Context, orderID domain.OrderID) PayOrderResult {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return PayOrderResult{OrderID: orderID, Success: false, Error: err}
	}

	if err := order.Pay(); err != nil {
		return PayOrderResult{OrderID: orderID, Success: false, Error: err}
	}

	if err := uc.paymentGW.Charge(ctx, orderID, order.Total()); err != nil {
		return PayOrderResult{OrderID: orderID, Success: false, Error: err}
	}

	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return PayOrderResult{OrderID: orderID, Success: false, Error: err}
	}

	return PayOrderResult{OrderID: orderID, Success: true}
}
