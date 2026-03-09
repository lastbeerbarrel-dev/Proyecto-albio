package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ItemMetadata represents static data about an item.
type ResourceAmount struct {
	ResourceID string `json:"resource_id"`
	Count      int    `json:"count"`
}

type ItemMetadata struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Tier           int              `json:"tier"`
	Category       string           `json:"category"` // e.g. "Weapon", "Armor", "Consumable", "Resource"
	SubCategory    string           `json:"sub_category,omitempty"`
	BaseIP         int              `json:"base_ip"`
	ReferencePrice int              `json:"reference_price"` // Default/Legacy
	HighVolatility bool             `json:"high_volatility,omitempty"`
	Recipe         []ResourceAmount `json:"recipe,omitempty"`
}

var itemDB = map[string]ItemMetadata{
	"T4_BAG":              {ID: "T4_BAG", Name: "Adept's Bag", Tier: 4, BaseIP: 700, ReferencePrice: 3000},
	"T4_MAIN_SWORD":       {ID: "T4_MAIN_SWORD", Name: "Adept's Broadsword", Tier: 4, BaseIP: 700, ReferencePrice: 4500},
	"T8_MAIN_ROCK_AVALON": {ID: "T8_MAIN_ROCK_AVALON", Name: "Avalonian Rock", Tier: 8, BaseIP: 0, ReferencePrice: 1100000},
	"T7_POTION_HEAL":      {ID: "T7_POTION_HEAL", Name: "Major Healing Potion", Tier: 7, BaseIP: 0, ReferencePrice: 5800},
	"T8_HEAD_PLATE":       {ID: "T8_HEAD_PLATE", Name: "Elder's Guardian Helmet", Tier: 8, BaseIP: 1100, ReferencePrice: 65000},
}

const dbFile = "metadata.json"

func loadMetadata() {
	file, err := os.Open(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No metadata file found, using defaults.")
			return
		}
		log.Printf("Error opening metadata file: %v", err)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&itemDB); err != nil {
		log.Printf("Error decoding metadata file: %v", err)
		return
	}
	log.Printf("Loaded %d items from storage.", len(itemDB))
}

func saveMetadata() {
	file, err := os.Create(dbFile)
	if err != nil {
		log.Printf("Error creating metadata file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(itemDB); err != nil {
		log.Printf("Error encoding metadata file: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	loadMetadata()
	// Save immediately to ensure file exists with defaults if it was missing
	saveMetadata()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(itemDB)
			return
		}

		if r.Method == "POST" {
			var newItems map[string]ItemMetadata
			if err := json.NewDecoder(r.Body).Decode(&newItems); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Merge/Overwrite
			for k, v := range newItems {
				itemDB[k] = v
			}
			saveMetadata()

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "updated", "count": fmt.Sprintf("%d", len(newItems))})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := r.URL.Path[len("/item/"):]
		item, exists := itemDB[id]
		if !exists {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metadata Service Operational"))
	})

	log.Printf("Metadata Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
