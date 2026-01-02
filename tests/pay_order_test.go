package tests

import (
	"context"
	"testing"

	"laba_7/application"
	"laba_7/domain"
	"laba_7/infrastructure/gateways"
	"laba_7/infrastructure/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayOrder_Success(t *testing.T) {
	ctx := context.Background()

	repo := repositories.NewInMemoryOrderRepository()
	gw := gateways.NewFakePaymentGateway()
	uc := application.NewPayOrderUseCase(repo, gw)

	orderID := domain.OrderID("order-123")
	order, err := domain.NewOrder(orderID, []domain.OrderLine{
		{ProductID: "p1", Quantity: 2, Price: domain.NewMoney(5000, "RUB")},
		{ProductID: "p2", Quantity: 1, Price: domain.NewMoney(10000, "RUB")},
	})
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, order))

	result := uc.Execute(ctx, orderID)
	assert.True(t, result.Success)
	assert.NoError(t, result.Error)
}

func TestPayOrder_EmptyOrder(t *testing.T) {
	_, err := domain.NewOrder("empty", nil)
	assert.ErrorIs(t, err, domain.ErrEmptyOrder)
}

func TestPayOrder_AlreadyPaid(t *testing.T) {
	ctx := context.Background()

	repo := repositories.NewInMemoryOrderRepository()
	gw := gateways.NewFakePaymentGateway()
	uc := application.NewPayOrderUseCase(repo, gw)

	orderID := domain.OrderID("paid-order")
	order, err := domain.NewOrder(orderID, []domain.OrderLine{
		{ProductID: "p1", Quantity: 1, Price: domain.NewMoney(1000, "RUB")},
	})
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, order))

	r1 := uc.Execute(ctx, orderID)
	require.True(t, r1.Success)

	r2 := uc.Execute(ctx, orderID)
	assert.False(t, r2.Success)
	assert.ErrorIs(t, r2.Error, domain.ErrAlreadyPaid)
}

func TestOrder_CannotModifyAfterPaid(t *testing.T) {
	order, err := domain.NewOrder("test", []domain.OrderLine{
		{ProductID: "p1", Quantity: 1, Price: domain.NewMoney(1000, "RUB")},
	})
	require.NoError(t, err)

	require.NoError(t, order.Pay())

	err = order.AddLine(domain.OrderLine{ProductID: "p2", Quantity: 1, Price: domain.NewMoney(2000, "RUB")})
	assert.ErrorIs(t, err, domain.ErrCannotModifyPaid)
}

func TestOrder_TotalValidation(t *testing.T) {
	order, err := domain.NewOrder("test", []domain.OrderLine{
		{ProductID: "p1", Quantity: 2, Price: domain.NewMoney(5000, "RUB")},  // 10000
		{ProductID: "p2", Quantity: 1, Price: domain.NewMoney(10000, "RUB")}, // 10000
	})
	require.NoError(t, err)

	assert.Equal(t, int64(20000), order.Total().Amount())
	assert.Equal(t, "RUB", order.Total().Currency())
}
