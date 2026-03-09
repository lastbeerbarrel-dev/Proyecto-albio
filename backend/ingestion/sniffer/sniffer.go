package sniffer

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"albion/common/models"
)

var mockItems = []struct {
	ID    string
	Price int
}{
	{"T4_BAG", 2500},            // High ROI
	{"T4_MAIN_SWORD", 3500},     // Medium ROI
	{"T8_MAIN_ROCK_AVALON", 900000}, // High Value
	{"T7_POTION_HEAL", 4000},    // Consumable
	{"T8_HEAD_PLATE", 60000},    // Gear
}

// Sniffer defines the interface for capturing network traffic.
type Sniffer interface {
	Start() error
	Stop()
	Updates() <-chan models.MarketOrder
}

// MockSniffer is a mock implementation of the Sniffer interface for environments without CGO/Npcap.
type MockSniffer struct {
	updates   chan models.MarketOrder
	stop      chan struct{}
	mockItems []struct {
		ID    string
		Price int
	}
	mutex sync.Mutex
}

// NewMockSniffer creates a new MockSniffer.
func NewMockSniffer() *MockSniffer {
	return &MockSniffer{
		updates: make(chan models.MarketOrder, 100),
		stop:    make(chan struct{}),
		mockItems: []struct {
			ID    string
			Price int
		}{
			{"T4_BAG", 2500},
			{"T4_MAIN_SWORD", 3500},
		},
	}
}

func (s *MockSniffer) UpdateMockItems(items []struct {
	ID    string
	Price int
}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.mockItems = items
}

func (s *MockSniffer) Start() error {
	log.Println("Mock Sniffer started. Generating fake market data...")
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				s.mutex.Lock()
				if len(s.mockItems) == 0 {
					s.mutex.Unlock()
					continue
				}
				// Pick random item
				item := s.mockItems[rand.Intn(len(s.mockItems))]
				s.mutex.Unlock()
				
				// Buy price: 60% to 90% of ref price to guarantee some profit
				buyFactor := 0.6 + (rand.Float64() * 0.3)
				price := int(float64(item.Price) * buyFactor)

				// Randomly add enchantment (20% chance)
				itemID := item.ID
				if rand.Float64() < 0.2 {
					enchant := rand.Intn(4) + 1
					itemID = itemID + fmt.Sprintf("@%d", enchant)
				}

				regions := []string{"Americas", "Europe", "Asia"}
				region := regions[rand.Intn(len(regions))]

				locations := []int{3005, 0007, 1002, 2001} // Caerleon, Swamps, etc.
				loc := locations[rand.Intn(len(locations))]

				// Emit a fake market order
				s.updates <- models.MarketOrder{
					OrderID:     time.Now().UnixNano(),
					ItemID:      itemID,
					LocationID:  loc,
					Quality:     rand.Intn(3) + 1,
					UnitPrice:   price,
					Amount:      rand.Intn(10) + 1,
					AuctionType: "offer",
					Expires:     time.Now().Add(24 * time.Hour),
					UpdatedAt:   time.Now(),
					Source:      "local_mock",
					Confidence:  "high",
					Region:      region,
				}
			}
		}
	}()
	return nil
}

func (s *MockSniffer) Stop() {
	close(s.stop)
	close(s.updates)
}

func (s *MockSniffer) Updates() <-chan models.MarketOrder {
	return s.updates
}
