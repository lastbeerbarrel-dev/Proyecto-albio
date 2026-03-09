package aodp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"albion/common/models"
)

var HostMap = map[string]string{
	"Americas": "www.albion-online-data.com",
	"Europe":   "europe.albion-online-data.com",
	"Asia":     "east.albion-online-data.com",
}

const APIPath = "/api/v2/stats/prices"

// City to ID Mapping
var CityToID = map[string]int{
	"Caerleon":      3005,
	"Black Market":  3003,
	"Bridgewatch":   1002,
	"Lymhurst":      1006,
	"Martlock":      3008,
	"Fort Sterling": 4002,
	"Thetford":      0007,
}

// Client is the client for the Albion Online Data Project API.
type Client struct {
	HTTPClient *http.Client
}

// NewClient creates a new AODP client.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchPrices gets current prices for a list of items at specific locations.
// PriceResponse represents the JSON response from AODP.
type PriceResponse struct {
	ItemID           string `json:"item_id"`
	City             string `json:"city"`
	Quality          int    `json:"quality"`
	SellPriceMin     int    `json:"sell_price_min"`
	SellPriceMinDate string `json:"sell_price_min_date"`
	BuyPriceMax      int    `json:"buy_price_max"`
	BuyPriceMaxDate  string `json:"buy_price_max_date"`
}

// FetchPrices gets current prices for a list of items at specific locations.
func (c *Client) FetchPrices(itemIDs []string, locations []string, region string) ([]models.MarketOrder, error) {
	if len(itemIDs) == 0 {
		return nil, nil
	}

	host, ok := HostMap[region]
	if !ok {
		host = HostMap["Americas"] // Fallback
	}

	itemsStr := strings.Join(itemIDs, ",")
	locsStr := url.QueryEscape(strings.Join(locations, ","))
	
	fetchURL := fmt.Sprintf("https://%s%s/%s?locations=%s", host, APIPath, itemsStr, locsStr)
    
	resp, err := c.HTTPClient.Get(fetchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AODP API returned status: %d", resp.StatusCode)
	}

	var prices []PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}

	var orders []models.MarketOrder
	for _, p := range prices {
		locID := CityToID[p.City]
		
		sellTime := parseTime(p.SellPriceMinDate)
		buyTime := parseTime(p.BuyPriceMaxDate)

		// Create orders for each quality found
		if p.SellPriceMin > 0 {
			orders = append(orders, models.MarketOrder{
				OrderID:     time.Now().UnixNano(),
				ItemID:      p.ItemID, // Keep original ID (e.g. T4_BAG)
				UnitPrice:   p.SellPriceMin,
				LocationID:  locID,
				Quality:     p.Quality,
				AuctionType: "offer",
				UpdatedAt:   sellTime,
				Source:      "aodp",
				Region:      region,
			})
		}
		if p.BuyPriceMax > 0 {
			orders = append(orders, models.MarketOrder{
				OrderID:     time.Now().UnixNano(),
				ItemID:      p.ItemID,
				UnitPrice:   p.BuyPriceMax,
				LocationID:  locID,
				Quality:     p.Quality,
				AuctionType: "request",
				UpdatedAt:   buyTime,
				Source:      "aodp",
				Region:      region,
			})
		}
	}

	return orders, nil
}

func parseTime(s string) time.Time {
	if s == "" || strings.HasPrefix(s, "0001") {
		return time.Now().Add(-24 * time.Hour) // Treat as stale
	}
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Now()
		}
	}
	return t
}
