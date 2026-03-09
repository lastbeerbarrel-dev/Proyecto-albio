package models

import "time"

// Item represents the static data of an item in Albion Online.
type Item struct {
	ItemID            string  `json:"ItemID"`
	ItemValue         int     `json:"ItemValue"` // Base value mainly for refining fees
	Tier              int     `json:"Tier"`
	EnchantmentLevel  int     `json:"EnchantmentLevel"`
	Quality           int     `json:"Quality"`
	MasteryModifier   float64 `json:"MasteryModifier"`
	BasePowerCapacity float64 `json:"BasePowerCapacity"` // For nutrition/crafting capacity
}

// MarketOrder represents a market order (buy or sell) from the data stream.
type MarketOrder struct {
	OrderID     int64     `json:"OrderID"`
	ItemID      string    `json:"ItemID"`
	LocationID  int       `json:"LocationID"`
	Quality     int       `json:"Quality"`
	UnitPrice   int       `json:"UnitPrice"` // Silver
	Amount      int       `json:"Amount"`
	AuctionType string    `json:"AuctionType"` // "offer" (sell) or "request" (buy)
	Expires     time.Time `json:"Expires"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	Source      string    `json:"Source"`     // "alo" (Albion Online Data Project) or "local" (Sniffer)
	Confidence  string    `json:"Confidence"` // "high", "medium", "low"
	Region      string    `json:"Region"`     // "Americas", "Europe", "Asia"
}

// TradeRoute represents a profitable path between two cities.
type TradeRoute struct {
	ItemID          string    `json:"item_id"`
	ItemName        string    `json:"item_name"`
	FromCity        string    `json:"from_city"`
	ToCity          string    `json:"to_city"`
	BuyPrice        int       `json:"buy_price"`
	SellPrice       int       `json:"sell_price"`
	Profit          float64   `json:"profit"`
	ROI             float64   `json:"roi"`
	Category        string    `json:"category"`
	SubCategory     string    `json:"sub_category"`
	Quality         int       `json:"quality"`
	ConfidenceScore int       `json:"confidence_score"` // Legacy mapping to ExecutionScore
	DataQuality     int       `json:"data_quality"`
	ExecutionScore  int       `json:"execution_score"`
	RiskProfile     string    `json:"risk_profile"`
	IsFiltered      bool      `json:"is_filtered"`
	HighVolatility  bool      `json:"high_volatility"`
	UpdatedAt       time.Time `json:"updated_at"`
}
