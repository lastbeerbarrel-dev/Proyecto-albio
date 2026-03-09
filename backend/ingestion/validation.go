package main

import (
	"albion/common/models"
	"fmt"
)

// ValidateOrder checks if the market order data is sane.
func ValidateOrder(order models.MarketOrder) error {
	if order.ItemID == "" {
		return fmt.Errorf("missing ItemID")
	}
	if order.UnitPrice <= 0 {
		return fmt.Errorf("invalid UnitPrice: %d", order.UnitPrice)
	}
	if order.Amount <= 0 {
		return fmt.Errorf("invalid Amount: %d", order.Amount)
	}
	// Add more checks (e.g. realistic max price)
	return nil
}
