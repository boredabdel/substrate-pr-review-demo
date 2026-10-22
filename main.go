package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Order struct {
	ID       string  `json:"id"`
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
	Status   string  `json:"status"`
}

type OrderStore struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[string]Order),
	}
}

func (s *OrderStore) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	o.Status = "CREATED"
	s.mu.Lock()
	s.orders[o.ID] = o
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(o)
}

func main() {
	store := NewOrderStore()
	http.HandleFunc("/orders", store.handleCreateOrder)
	fmt.Println("Order service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
