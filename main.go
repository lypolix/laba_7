package main

import (
	"context"
	"fmt"
	"log"

	"laba_7/application"
	"laba_7/domain"
	"laba_7/infrastructure/gateways"
	"laba_7/infrastructure/repositories"
)

func main() {
	ctx := context.Background()

	repo := repositories.NewInMemoryOrderRepository()
	gw := gateways.NewFakePaymentGateway()

	uc := application.NewPayOrderUseCase(repo, gw)

	orderID := domain.OrderID("order-123")
	result := uc.Execute(ctx, orderID)

	if result.Success {
		fmt.Printf("Order %s paid successfully\n", orderID)
	} else {
		log.Printf("Failed to pay order %s: %v", orderID, result.Error)
	}
}
