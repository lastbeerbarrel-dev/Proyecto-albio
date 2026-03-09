package repository

import (
	"errors"
	"strings"

	"albion/common/models"
)

// ItemRepository defines the interface for retrieving item data.
type ItemRepository interface {
	GetItem(itemID string) (models.Item, error)
}

// InMemoryRepository is a mock implementation for MVP/No-DB environments.
type InMemoryRepository struct {
	items map[string]models.Item
}

// NewInMemoryRepository creates a new repository with pre-populated data.
func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{
		items: make(map[string]models.Item),
	}
	repo.seed()
	return repo
}

func (r *InMemoryRepository) seed() {
	// Sample Data based on spec
	r.items["T8_MAIN_ROCK_AVALON"] = models.Item{
		ItemID:            "T8_MAIN_ROCK_AVALON",
		ItemValue:         5120,
		Tier:              8,
		EnchantmentLevel:  0,
		Quality:           1,
		MasteryModifier:   0.20,
		BasePowerCapacity: 144,
	}
	r.items["T4_BAG"] = models.Item{
		ItemID:            "T4_BAG",
		ItemValue:         16,
		Tier:              4,
		EnchantmentLevel:  0,
		Quality:           1,
		MasteryModifier:   0,
		BasePowerCapacity: 0,
	}
}

func (r *InMemoryRepository) GetItem(itemID string) (models.Item, error) {
	// Simple lookup, handling enchantment/quality suffixes might be needed in real app
	// e.g. T4_BAG@1
	baseID := strings.Split(itemID, "@")[0] // Simplistic stripping
	
	if item, ok := r.items[baseID]; ok {
		return item, nil
	}
	return models.Item{}, errors.New("item not found")
}
