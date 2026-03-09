package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"albion/common/models"
	"albion/ingestion/aodp"
	"albion/ingestion/sniffer"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Calculation Service URL (Hardcoded for Dev)
	calcServiceURL := "http://localhost:8081/ingest"

	// Initialize AODP Client
	aodpClient := aodp.NewClient()
	_ = aodpClient

	// Initialize Sniffer
	var sniff sniffer.Sniffer

	if os.Getenv("USE_REAL_SNIFFER") == "true" {
		var err error
		sniff, err = sniffer.NewRealSniffer("")
		if err != nil {
			log.Fatalf("Failed to initialize Real Sniffer: %v", err)
		}
		if err := sniff.Start(); err != nil {
			log.Fatalf("Failed to start sniffer: %v", err)
		}
		defer sniff.Stop()
	} else {
		log.Println("Simulation (Mock Sniffer) disabled. Using AODP Real Market Data only.")
	}

	// Always enforce Real Market Data/Sniffer
	log.Println("Real Data Polling (AODP) and Sniffer (if enabled) active.")

	// Start processing Sniffer updates (if one is active)
	if sniff != nil {
		go func() {
			client := &http.Client{}
			for order := range sniff.Updates() {
				// Validate Data
				if err := ValidateOrder(order); err != nil {
					log.Printf("Skipping invalid order: %v", err)
					continue
				}

				log.Printf("Received local order: %s - %d", order.ItemID, order.UnitPrice)
				pushOrder(client, calcServiceURL, order)
			}
		}()
	}

	// Real Data Polling Loop (AODP)
	go func() {
		client := &http.Client{}
		cities := []string{"Caerleon", "Black Market", "Bridgewatch", "Lymhurst", "Martlock", "Fort Sterling", "Thetford"}
		regions := []string{"Americas", "Europe", "Asia"}

		for {
			// 1. Get current item list from Metadata
			res, err := http.Get("http://localhost:8082/items")
			if err == nil {
				var items map[string]interface{}
				if json.NewDecoder(res.Body).Decode(&items) == nil {
					var itemIDs []string
					for k := range items {
						itemIDs = append(itemIDs, k)
					}

					// 2. Fetch from AODP for each region in chunks to avoid URL length issues
					if len(itemIDs) > 0 {
						log.Printf("Polling AODP for %d items across %d regions", len(itemIDs), len(regions))
						batchSize := 20
						for _, reg := range regions {
							for i := 0; i < len(itemIDs); i += batchSize {
								end := i + batchSize
								if end > len(itemIDs) {
									end = len(itemIDs)
								}
								chunk := itemIDs[i:end]

								// Fetch quality 1 explicitly to ensure we get liquid data from AODP
								orders, err := aodpClient.FetchPrices(chunk, cities, reg)
								if err == nil {
									for _, o := range orders {
										pushOrder(client, calcServiceURL, o)
									}
									log.Printf("Fetched %d real orders for %s (chunk %d-%d) from AODP", len(orders), reg, i, end)
								} else {
									log.Printf("Error fetching AODP data for %s: %v", reg, err)
								}
								time.Sleep(5 * time.Second) // Throttle to avoid rate limits
							}
						}
					}
				}
				res.Body.Close()
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ingestion Service Operational"))
	})

	log.Printf("Ingestion Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func pushOrder(client *http.Client, url string, order models.MarketOrder) {
	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshaling order: %v", err)
		return
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error pushing order to calculation service: %v", err)
		return
	}
	resp.Body.Close()
}
