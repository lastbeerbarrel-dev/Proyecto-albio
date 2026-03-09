package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ItemMetadata struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	Tier           int                      `json:"tier"`
	Category       string                   `json:"category"`
	SubCategory    string                   `json:"sub_category"`
	BaseIP         int                      `json:"base_ip"`
	ReferencePrice int                      `json:"reference_price"`
	Prices         map[string]int           `json:"prices,omitempty"`
	Recipe         []map[string]interface{} `json:"recipe,omitempty"`
}

type AODPResponse struct {
	ItemID           string `json:"item_id"`
	City             string `json:"city"`
	Quality          int    `json:"quality"`
	SellPriceMin     int    `json:"sell_price_min"`
	SellPriceMinDate string `json:"sell_price_min_date"`
}

var regionHosts = map[string]string{
	"Americas": "www.albion-online-data.com",
	"Europe":   "europe.albion-online-data.com",
	"Asia":     "east.albion-online-data.com",
}

func main() {
	log.Println("Starting reference price population from AODP...")

	// Read current metadata
	data, err := os.ReadFile("backend/metadata/metadata.json")
	if err != nil {
		log.Fatalf("Failed to read metadata.json: %v", err)
	}

	var items map[string]ItemMetadata
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatalf("Failed to parse metadata.json: %v", err)
	}

	log.Printf("Loaded %d items from metadata", len(items))

	// Collect all item IDs
	var itemIDs []string
	for id := range items {
		itemIDs = append(itemIDs, id)
	}

	// Process in batches
	batchSize := 20
	regions := []string{"Americas", "Europe", "Asia"}
	client := &http.Client{Timeout: 30 * time.Second}

	updatedCount := 0
	for i := 0; i < len(itemIDs); i += batchSize {
		end := i + batchSize
		if end > len(itemIDs) {
			end = len(itemIDs)
		}
		batch := itemIDs[i:end]

		log.Printf("Processing batch %d-%d of %d items...", i+1, end, len(itemIDs))

		// Fetch prices for each region
		for _, region := range regions {
			prices, err := fetchPricesFromAODP(client, batch, region)
			if err != nil {
				log.Printf("Error fetching prices for %s: %v", region, err)
				continue
			}

			// Update items with prices
			for itemID, avgPrice := range prices {
				if item, exists := items[itemID]; exists {
					if item.Prices == nil {
						item.Prices = make(map[string]int)
					}
					item.Prices[region] = avgPrice
					items[itemID] = item
					updatedCount++
				}
			}
		}

		// Throttle to avoid rate limiting
		if i+batchSize < len(itemIDs) {
			log.Printf("Waiting 5 seconds before next batch...")
			time.Sleep(5 * time.Second)
		}
	}

	log.Printf("Updated prices for %d item-region combinations", updatedCount)

	// Save updated metadata
	output, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal updated metadata: %v", err)
	}

	if err := os.WriteFile("backend/metadata/metadata.json", output, 0644); err != nil {
		log.Fatalf("Failed to write updated metadata: %v", err)
	}

	log.Println("Successfully updated metadata.json with reference prices!")
	log.Printf("Total items processed: %d", len(items))
	log.Printf("Total price updates: %d", updatedCount)
}

func fetchPricesFromAODP(client *http.Client, itemIDs []string, region string) (map[string]int, error) {
	host := regionHosts[region]
	itemsStr := strings.Join(itemIDs, ",")
	locations := url.QueryEscape("Caerleon,Black Market,Bridgewatch,Lymhurst,Martlock,Fort Sterling,Thetford")

	fetchURL := fmt.Sprintf("https://%s/api/v2/stats/prices/%s?locations=%s", host, itemsStr, locations)

	resp, err := client.Get(fetchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AODP returned status %d", resp.StatusCode)
	}

	var aodpData []AODPResponse
	if err := json.NewDecoder(resp.Body).Decode(&aodpData); err != nil {
		return nil, err
	}

	// Calculate average prices per item
	priceMap := make(map[string][]int)
	for _, entry := range aodpData {
		if entry.SellPriceMin > 0 {
			priceMap[entry.ItemID] = append(priceMap[entry.ItemID], entry.SellPriceMin)
		}
	}

	// Calculate averages
	avgPrices := make(map[string]int)
	for itemID, prices := range priceMap {
		if len(prices) > 0 {
			sum := 0
			for _, p := range prices {
				sum += p
			}
			avgPrices[itemID] = sum / len(prices)
		}
	}

	return avgPrices, nil
}
