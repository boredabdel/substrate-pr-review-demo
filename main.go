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

// ApplyDiscount applies a validated percentage discount to an existing order safely.
func (s *OrderStore) handleApplyDiscount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	var req struct {
		DiscountPercent float64 `json:"discountPercent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}
	if req.DiscountPercent <= 0 || req.DiscountPercent > 100 {
		http.Error(w, "discountPercent must be between 0 and 100", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	order, exists := s.orders[orderID]
	if !exists {
		s.mu.Unlock()
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	order.Amount = order.Amount * (1 - req.DiscountPercent/100)
	s.orders[orderID] = order
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

func main() {
	store := NewOrderStore()
	http.HandleFunc("/orders", store.handleCreateOrder)
	http.HandleFunc("/orders/discount", store.handleApplyDiscount)
	fmt.Println("Order service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
