package messaging

import (
	"albion/common/models"
)

// Publisher defines the interface for publishing market data.
type Publisher interface {
	PublishOrder(order models.MarketOrder) error
}

// Subscriber defines the interface for consuming market data.
type Subscriber interface {
	SubscribeOrders() (<-chan models.MarketOrder, error)
}

// MockMessaging is an in-memory implementation of Publisher and Subscriber.
type MockMessaging struct {
	orders chan models.MarketOrder
}

// NewMockMessaging creates a new MockMessaging instance.
func NewMockMessaging() *MockMessaging {
	return &MockMessaging{
		orders: make(chan models.MarketOrder, 1000),
	}
}

func (m *MockMessaging) PublishOrder(order models.MarketOrder) error {
	select {
	case m.orders <- order:
		return nil
	default:
		// Drop if full
		return nil
	}
}

func (m *MockMessaging) SubscribeOrders() (<-chan models.MarketOrder, error) {
	return m.orders, nil
}
