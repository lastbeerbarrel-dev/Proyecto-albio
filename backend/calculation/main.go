package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"albion/calculation/engine"
	"albion/common/models"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan engine.Opportunity
	mutex     sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan engine.Opportunity),
		clients:   make(map[*websocket.Conn]bool),
	}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	hub  = newHub()
	calc = engine.NewCalculator()

	// Stats tracking
	statsMutex  sync.Mutex
	totalProfit float64
	activeFlips int
	recentFlips = make(map[string]time.Time)
)

func main() {
	// Sync reference prices and categories from Metadata Service periodically
	go func() {
		for {
			res, err := http.Get("http://localhost:8082/items")
			if err == nil {
				var items map[string]struct {
					ReferencePrice int            `json:"reference_price"`
					Prices         map[string]int `json:"prices"`
					Category       string         `json:"category"`
					SubCategory    string         `json:"sub_category"`
					Name           string         `json:"name"`
					HighVolatility bool           `json:"high_volatility"`
				}
				if json.NewDecoder(res.Body).Decode(&items) == nil {
					regionalPrices := make(map[string]map[string]int)
					categories := make(map[string]string)
					subCats := make(map[string]string)
					names := make(map[string]string)
					highVol := make(map[string]bool)
					for k, v := range items {
						categories[k] = v.Category
						subCats[k] = v.SubCategory
						names[k] = v.Name
						highVol[k] = v.HighVolatility
						if v.Prices != nil {
							regionalPrices[k] = v.Prices
						} else {
							regionalPrices[k] = map[string]int{
								"Americas": v.ReferencePrice,
								"Europe":   v.ReferencePrice,
								"Asia":     v.ReferencePrice,
							}
						}
					}
					calc.UpdateRegionalPrices(regionalPrices)
					calc.UpdateMetadata(categories, subCats, names, highVol)
					log.Printf("Loaded %d regional reference prices and metadata", len(regionalPrices))
				}
				res.Body.Close()
			}
			time.Sleep(30 * time.Second)
		}
	}()

	// Stats Cleanup Loop (Remove stale flips from active count)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			statsMutex.Lock()
			now := time.Now()
			for id, t := range recentFlips {
				if now.Sub(t) > 10*time.Minute {
					delete(recentFlips, id)
				}
			}
			activeFlips = len(recentFlips)
			statsMutex.Unlock()
		}
	}()

	// Stats Endpoint
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		statsMutex.Lock()
		defer statsMutex.Unlock()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_profit": totalProfit,
			"active_flips": activeFlips,
		})
	})

	// Trade Routes Endpoint
	http.HandleFunc("/calculate/routes", func(w http.ResponseWriter, r *http.Request) {
		region := r.URL.Query().Get("region")
		if region == "" {
			region = "Americas"
		}
		routes := calc.GetTradeRoutes(region)
		json.NewEncoder(w).Encode(routes)
	})

	// Premium Toggle Handler
	http.HandleFunc("/calculate/premium", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var req struct {
				Premium bool `json:"premium"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			calc.SetPremium(req.Premium)
			w.WriteHeader(http.StatusOK)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"premium": calc.IsPremium})
	})

	// Ingestion Handler
	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		var order models.MarketOrder
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			log.Printf("Error decoding ingest order: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Received order for %s at %d (Price: %d, Q: %d)", order.ItemID, order.LocationID, order.UnitPrice, order.Quality)
		opp := calc.CalculateFlip(order)
		hub.broadcast <- opp
		w.WriteHeader(http.StatusOK)
	})

	// WebSocket Handler
	http.HandleFunc("/ws", handleWebSocket)

	// Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go hub.run()

	// CORS Middleware
	withCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	log.Println("Calculation Engine starting on port 8081")
	log.Fatal(http.ListenAndServe(":8081", withCORS(http.DefaultServeMux)))
}

func (h *Hub) run() {
	for {
		opp := <-h.broadcast

		// Update Stats
		if opp.Profit > 0 {
			statsMutex.Lock()
			totalProfit += opp.Profit
			recentFlips[opp.OriginalOrder.ItemID] = time.Now()
			activeFlips = len(recentFlips)
			statsMutex.Unlock()
		}

		h.mutex.Lock()
		for client := range h.clients {
			err := client.WriteJSON(opp)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mutex.Unlock()
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.mutex.Lock()
	hub.clients[conn] = true
	hub.mutex.Unlock()
}
