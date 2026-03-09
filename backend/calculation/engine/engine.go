package engine

import (
	"albion/common/models"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Opportunity represents a calculated market flip.
type Opportunity struct {
	OriginalOrder  models.MarketOrder `json:"OriginalOrder"`
	Profit         float64            `json:"Profit"`
	ROI            float64            `json:"ROI"`
	Confidence     string             `json:"Confidence"` // Legacy
	DataQuality    int                `json:"DataQuality"`
	ExecutionScore int                `json:"ExecutionScore"`
	RiskProfile    string             `json:"RiskProfile"`
	Category       string             `json:"Category"`
	SubCategory    string             `json:"SubCategory"`
	ItemName       string             `json:"ItemName"`
	BuyLocation    string             `json:"BuyLocation"`
	SellLocation   string             `json:"SellLocation"`
	SellPrice      float64            `json:"SellPrice"`
	HighVolatility bool               `json:"HighVolatility"`
}

// Calculator handles complex game logic calculations.
type Calculator struct {
	CraftingFeeMultiplier float64
	MarketTaxRate         float64                                            // Current rate
	SetupFeeRate          float64                                            // e.g. 0.025 (2.5%)
	RegionalRefPrices     map[string]map[string]int                          // ItemID -> Region -> Price
	Categories            map[string]string                                  // ItemID -> Category
	SubCategories         map[string]string                                  // ItemID -> SubCategory
	Names                 map[string]string                                  // ItemID -> DisplayName
	CurrentBuyPrices      map[string]map[string]map[int]int                  // Best Buy Order (Request)
	CurrentSellPrices     map[string]map[string]map[int]int                  // Best Sell Offer (Offer)
	CurrentPriceTimes     map[string]map[string]map[int]time.Time            // ItemID -> Region -> LocID -> Time
	OrderBuffer           map[string]map[string]map[int][]models.MarketOrder // History for Buy Orders (Outlier detection)
	HighVolatility        map[string]bool                                    // ItemID -> HighVolatility
	LocationNames         map[int]string
	IsPremium             bool
}

func NewCalculator() *Calculator {
	return &Calculator{
		CraftingFeeMultiplier: 1.0,
		MarketTaxRate:         0.04, // Default Premium
		SetupFeeRate:          0.025,
		RegionalRefPrices:     make(map[string]map[string]int),
		Categories:            make(map[string]string),
		SubCategories:         make(map[string]string),
		Names:                 make(map[string]string),
		CurrentBuyPrices:      make(map[string]map[string]map[int]int),
		CurrentSellPrices:     make(map[string]map[string]map[int]int),
		CurrentPriceTimes:     make(map[string]map[string]map[int]time.Time),
		OrderBuffer:           make(map[string]map[string]map[int][]models.MarketOrder),
		HighVolatility:        make(map[string]bool),
		LocationNames: map[int]string{
			3005: "Caerleon",
			1002: "Bridgewatch",
			1006: "Lymhurst",
			3008: "Martlock",
			4002: "Fort Sterling",
			0007: "Thetford",
			3003: "Mercado Negro",
		},
		IsPremium: true,
	}
}

func (c *Calculator) UpdateRegionalPrices(prices map[string]map[string]int) {
	c.RegionalRefPrices = prices
}

func (c *Calculator) UpdateMetadata(categories map[string]string, subCats map[string]string, names map[string]string, highVol map[string]bool) {
	c.Categories = categories
	c.SubCategories = subCats
	c.Names = names
	c.HighVolatility = highVol
}

func (c *Calculator) SetPremium(premium bool) {
	c.IsPremium = premium
	if premium {
		c.MarketTaxRate = 0.04
	} else {
		c.MarketTaxRate = 0.08
	}
}

// CalculateDataQuality measures the statistical robustness of the data (0-100).
func (c *Calculator) CalculateDataQuality(itemID string, region string, locID int, quality int) int {
	compositeID := fmt.Sprintf("%s#%d", itemID, quality)
	orders, ok := c.OrderBuffer[compositeID][region][locID]
	if !ok || len(orders) == 0 {
		return 0
	}

	score := 0
	// 1. Quantity of orders (last 10)
	score += len(orders) * 4 // Max 40%

	// 2. Freshness
	lastUpdate := time.Now()
	if times, ok := c.CurrentPriceTimes[compositeID]; ok {
		if reg, ok := times[region]; ok {
			lastUpdate = reg[locID]
		}
	}
	age := time.Since(lastUpdate).Minutes()
	if age < 60 {
		score += 40
	} else if age < 360 {
		score += 20
	} else if age < 1440 {
		score += 10
	}

	// 3. Dispersion (simple check)
	if len(orders) >= 3 {
		score += 20
	}

	if score > 100 {
		score = 100
	}
	return score
}

// CalculateExecutionScore measures how profitable and executable a flip is (0-100).
func (c *Calculator) CalculateExecutionScore(profit float64, roi float64, itemID string, region string, locID int, quality int, category string) int {
	// 1. NetSpreadScore (0.45) - Profit and ROI
	netSpreadScore := 0.0
	if roi > 20 {
		netSpreadScore = 100
	} else if roi > 10 {
		netSpreadScore = 70
	} else if roi > 5 {
		netSpreadScore = 40
	}

	// 2. HistoricalBMAlignment (0.25)
	// Compare with reference price (proxy for P50)
	baseID := itemID
	if idx := strings.Index(itemID, "@"); idx != -1 {
		baseID = itemID[:idx]
	}
	bmAlignment := 0.0
	if regions, ok := c.RegionalRefPrices[baseID]; ok {
		if ref, ok := regions[region]; ok && ref > 0 {
			// If we sell at BM (loc 3003), check if sell price is near ref
			bmAlignment = 50 // Default neutral
		}
	}

	// 3. LiquidityTimeScore (0.20)
	liquidityScore := 0.0
	switch category {
	case "Consumable", "Resource":
		liquidityScore = 100
	case "Weapon", "Armor":
		liquidityScore = 60
	default:
		liquidityScore = 30
	}

	// 4. OrderSanityScore (0.10)
	sanity := 100.0 // Default healthy

	total := (netSpreadScore * 0.45) + (bmAlignment * 0.25) + (liquidityScore * 0.20) + (sanity * 0.10)

	// Black Market Special Rules (V2.0)
	if locID == 3003 {
		if profit >= 150000 {
			if total < 35 {
				total = 35
			}
		}
	}

	return int(total)
}

// CalculateRiskProfile provides a qualitative risk assessment.
func (c *Calculator) CalculateRiskProfile(dataQuality int, executionScore int, category string) string {
	if executionScore > 70 && dataQuality > 60 {
		return "Bajo"
	}
	if executionScore > 40 || dataQuality > 40 {
		return "Medio"
	}
	return "Alto"
}

// CalculateEffectivePrice uses IQR to filter outliers and returns a realistic price.
func (c *Calculator) CalculateEffectivePrice(itemID string, region string, locID int, quality int, category string) (float64, int) {
	compositeID := fmt.Sprintf("%s#%d", itemID, quality)
	orders, ok := c.OrderBuffer[compositeID][region][locID]
	if !ok || len(orders) == 0 {
		return 0, 0
	}

	// 1. Minimum Volume Filter
	minVol := 1
	if category == "Resource" || category == "Consumable" {
		minVol = 50
	}

	validOrders := []models.MarketOrder{}
	prices := []int{}

	for _, o := range orders {
		if o.Amount >= minVol {
			validOrders = append(validOrders, o)
			prices = append(prices, o.UnitPrice)
		}
	}

	if len(prices) == 0 {
		return 0, 0
	}

	// 2. Statistical Filter (IQR) if enough data
	if len(prices) >= 4 {
		sort.Ints(prices)
		q1 := prices[len(prices)/4]
		q3 := prices[len(prices)*3/4]
		iqr := q3 - q1
		lowerBound := float64(q1) - 1.5*float64(iqr)

		finalPrices := []int{}
		for _, p := range prices {
			if float64(p) >= lowerBound {
				finalPrices = append(finalPrices, p)
			}
		}
		if len(finalPrices) > 0 {
			sort.Ints(finalPrices)
			return float64(finalPrices[len(finalPrices)-1]), 100 // Return highest valid price
		}
	}

	// Fallback to Median if not enough for IQR
	sort.Ints(prices)
	median := prices[len(prices)/2]
	return float64(median), 50
}

// CalculateFlip determines potential profit and ROI for a market flip.
func (c *Calculator) CalculateFlip(order models.MarketOrder) Opportunity {
	compositeID := fmt.Sprintf("%s#%d", order.ItemID, order.Quality)

	// Update Buffer
	if _, ok := c.OrderBuffer[compositeID]; !ok {
		c.OrderBuffer[compositeID] = make(map[string]map[int][]models.MarketOrder)
	}
	if _, ok := c.OrderBuffer[compositeID][order.Region]; !ok {
		c.OrderBuffer[compositeID][order.Region] = make(map[int][]models.MarketOrder)
	}

	// Circular Buffer logic (keep last 10)
	buffer := c.OrderBuffer[compositeID][order.Region][order.LocationID]
	buffer = append(buffer, order)
	if len(buffer) > 10 {
		buffer = buffer[1:]
	}
	c.OrderBuffer[compositeID][order.Region][order.LocationID] = buffer

	// Update Current Prices (Differentiated)
	if order.AuctionType == "request" {
		if _, ok := c.CurrentBuyPrices[compositeID]; !ok {
			c.CurrentBuyPrices[compositeID] = make(map[string]map[int]int)
		}
		if _, ok := c.CurrentBuyPrices[compositeID][order.Region]; !ok {
			c.CurrentBuyPrices[compositeID][order.Region] = make(map[int]int)
		}
		c.CurrentBuyPrices[compositeID][order.Region][order.LocationID] = order.UnitPrice
	} else {
		if _, ok := c.CurrentSellPrices[compositeID]; !ok {
			c.CurrentSellPrices[compositeID] = make(map[string]map[int]int)
		}
		if _, ok := c.CurrentSellPrices[compositeID][order.Region]; !ok {
			c.CurrentSellPrices[compositeID][order.Region] = make(map[int]int)
		}
		// We track the LOWEST Sell Offer
		current := c.CurrentSellPrices[compositeID][order.Region][order.LocationID]
		if current == 0 || order.UnitPrice < current {
			c.CurrentSellPrices[compositeID][order.Region][order.LocationID] = order.UnitPrice
		}
	}

	// Track Time
	if _, ok := c.CurrentPriceTimes[compositeID]; !ok {
		c.CurrentPriceTimes[compositeID] = make(map[string]map[int]time.Time)
	}
	if _, ok := c.CurrentPriceTimes[compositeID][order.Region]; !ok {
		c.CurrentPriceTimes[compositeID][order.Region] = make(map[int]time.Time)
	}
	c.CurrentPriceTimes[compositeID][order.Region][order.LocationID] = order.UpdatedAt

	// Metadata Lookup
	baseID := order.ItemID
	enchant := ".0"
	if idx := strings.Index(order.ItemID, "@"); idx != -1 {
		baseID = order.ItemID[:idx]
		enchant = "." + order.ItemID[idx+1:]
	}

	category := c.Categories[baseID]
	displayName := c.Names[baseID]
	subCat := c.SubCategories[baseID]

	if displayName == "" {
		category = c.Categories[order.ItemID]
		displayName = c.Names[order.ItemID]
		subCat = c.SubCategories[order.ItemID]
	}

	if displayName == "" {
		displayName = order.ItemID
	}
	if enchant != ".0" && !strings.Contains(displayName, enchant) {
		displayName += " " + enchant
	}

	var potentialSellPrice float64

	// 1. Try to get real-time Buy Order price from Black Market (loc 3003)
	if bmPrices, ok := c.CurrentBuyPrices[compositeID]; ok {
		if bmRegion, ok := bmPrices[order.Region]; ok {
			if price := bmRegion[3003]; price > 0 {
				potentialSellPrice = float64(price)
			}
		}
	}

	// 2. Last fallback: if no potentialSellPrice, we can't calculate a profitable flip
	if potentialSellPrice == 0 {
		return Opportunity{OriginalOrder: order, Confidence: "low", Category: category, ItemName: displayName}
	}
	switch enchant {
	case ".1":
		potentialSellPrice *= 1.4
	case ".2":
		potentialSellPrice *= 1.8
	case ".3":
		potentialSellPrice *= 3.0
	case ".4":
		potentialSellPrice *= 6.0
	}

	isBlackMarketable := false
	switch category {
	case "Armor", "Weapon", "Accessory", "Offhand", "Tome":
		isBlackMarketable = true
	}

	buyLoc := "Mercado"
	if order.LocationID == 3005 {
		buyLoc = "Caerleon"
	}
	sellLoc := "Mercado Negro"

	taxRate := c.MarketTaxRate
	setupRate := c.SetupFeeRate

	// If we are selling to BM, we use the BM tax logic
	if isBlackMarketable {
		setupRate = 0.0 // Selling to a Buy Order (Market Request) ignores setup fee
		taxRate = 0.04  // Black Market tax is fixed at 4%? Or follows premium? Let's keep it safe.
	} else {
		sellLoc = "No apto para BM"
		// If not BM-able, we assume selling at a normal market
	}

	totalCosts := (potentialSellPrice * taxRate) + (potentialSellPrice * setupRate) + float64(order.UnitPrice)
	profit := potentialSellPrice - totalCosts
	roi := (profit / float64(order.UnitPrice)) * 100

	dataQual := c.CalculateDataQuality(order.ItemID, order.Region, 3003, order.Quality)
	execScore := c.CalculateExecutionScore(profit, roi, order.ItemID, order.Region, 3003, order.Quality, category)
	riskProf := c.CalculateRiskProfile(dataQual, execScore, category)

	highVol := c.HighVolatility[baseID]
	if highVol == false {
		highVol = c.HighVolatility[order.ItemID]
	}

	conf := "medium"
	if execScore > 70 {
		conf = "high"
	} else if execScore < 30 {
		conf = "low"
	}

	if !isBlackMarketable {
		profit = 0
		roi = 0
		conf = "none"
		execScore = 0
	}

	return Opportunity{
		OriginalOrder:  order,
		Profit:         profit,
		ROI:            roi,
		Confidence:     conf,
		DataQuality:    dataQual,
		ExecutionScore: execScore,
		RiskProfile:    riskProf,
		Category:       category,
		SubCategory:    subCat,
		ItemName:       displayName,
		BuyLocation:    buyLoc,
		SellLocation:   sellLoc,
		SellPrice:      potentialSellPrice,
		HighVolatility: highVol,
	}
}

// GetTradeRoutes identifies profitable paths between cities.
func (c *Calculator) GetTradeRoutes(region string) []models.TradeRoute {
	routes := []models.TradeRoute{}
	// We iterate through items that have SELL OFFERS (Purchase source)
	for compositeID, regions := range c.CurrentSellPrices {
		itemID := compositeID
		quality := 1
		parts := strings.Split(compositeID, "#")
		if len(parts) == 2 {
			itemID = parts[0]
			if q, err := strconv.Atoi(parts[1]); err == nil {
				quality = q
			}
		}

		sellCityPrices, ok := regions[region]
		if !ok {
			continue
		}

		baseID := itemID
		if idx := strings.Index(itemID, "@"); idx != -1 {
			baseID = itemID[:idx]
		}

		category := c.Categories[baseID]
		displayName := c.Names[baseID]
		subCat := c.SubCategories[baseID]

		if displayName == "" {
			category = c.Categories[itemID]
			displayName = c.Names[itemID]
			subCat = c.SubCategories[itemID]
		}
		if displayName == "" {
			displayName = itemID
		}

		isBM := false
		switch category {
		case "Armor", "Weapon", "Accessory", "Offhand", "Tome":
			isBM = true
		}

		// Purchase from Sell Offer at locA
		for locA, priceA := range sellCityPrices {
			if locA == 3003 || locA == 0 {
				continue
			}

			// We check ALL cities (including locA itself for city-flipping) for BUY ORDERS
			buyRegions, ok := c.CurrentBuyPrices[compositeID]
			if !ok {
				continue
			}
			buyCityPrices, ok := buyRegions[region]
			if !ok {
				continue
			}

			for locB, priceB := range buyCityPrices {
				if locB == 0 {
					continue
				}
				if locB == 3003 && !isBM {
					continue
				}

				pB := float64(priceB)
				pA := float64(priceA)

				// USE EFFECTIVE PRICE LOGIC (Statistical verification of the Buy Order)
				effPrice, _ := c.CalculateEffectivePrice(itemID, region, locB, quality, category)
				isFiltered := false
				if effPrice > 0 {
					if effPrice < pB {
						pB = effPrice
						isFiltered = true
					}
				}

				taxRate := c.MarketTaxRate
				setupRate := 0.0 // When selling to Buy Order, no setup fee
				if locB == 3003 {
					taxRate = 0.04 // Black Market tax is fixed? Or 0? Let's assume standard tax
				}

				profit := pB*(1.0-taxRate-setupRate) - pA

				highVol := c.HighVolatility[baseID]
				if highVol == false {
					highVol = c.HighVolatility[compositeID]
				}

				if profit > 1000 {
					roi := (profit / pA) * 100
					if roi > 3 {
						var updatedAt time.Time
						if times, ok := c.CurrentPriceTimes[compositeID]; ok {
							if regionTimes, ok := times[region]; ok {
								updatedAt = regionTimes[locB]
							}
						}

						dataQual := c.CalculateDataQuality(itemID, region, locB, quality)
						execScore := c.CalculateExecutionScore(profit, roi, itemID, region, locB, quality, category)
						riskProf := c.CalculateRiskProfile(dataQual, execScore, category)

						// One opportunity is shown if ExecutionScore >= 30
						if execScore >= 30 {
							routes = append(routes, models.TradeRoute{
								ItemID:          itemID,
								ItemName:        displayName,
								FromCity:        c.LocationNames[locA],
								ToCity:          c.LocationNames[locB],
								BuyPrice:        priceA,
								SellPrice:       int(pB),
								Profit:          profit,
								ROI:             roi,
								Category:        category,
								SubCategory:     subCat,
								Quality:         quality,
								ConfidenceScore: execScore,
								DataQuality:     dataQual,
								ExecutionScore:  execScore,
								RiskProfile:     riskProf,
								IsFiltered:      isFiltered,
								HighVolatility:  highVol,
								UpdatedAt:       updatedAt,
							})
						}
					}
				}
			}
		}
	}
	return routes
}
