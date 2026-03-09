package main

import (
	"encoding/json"
	"net/http"
)

// FetchReferencePrices gets prices from Metadata Service.
func FetchReferencePrices(url string) (map[string]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items map[string]struct {
		ReferencePrice int `json:"reference_price"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	
	prices := make(map[string]int)
	for id, item := range items {
		prices[id] = item.ReferencePrice
	}
	return prices, nil
}
