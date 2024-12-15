package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"realtime-dashboard-food-delivery/model"
	"realtime-dashboard-food-delivery/service"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	db           *sql.DB
	notification *service.Notification
}

func NewServer(db *sql.DB, notification *service.Notification) *Server {
	return &Server{db: db, notification: notification}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	r = r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")
	// Set CORS headers immediately, before any potential errors
	w.Header().Set("Access-Control-Allow-Origin", "*") // For testing, allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*") // For testing, allow all headers

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set default content type after CORS headers
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == "POST" && r.URL.Path == "/orders":
		s.handlerCreateOrder(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/poll/orders/"): // Added leading slash
		s.handlePollOrder(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *Server) handlerCreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := `
        INSERT INTO orders (restaurant_id, status)
        VALUES ($1, $2)
        RETURNING id, created_at`
	order.Status = "new"

	err = tx.QueryRow(query, order.RestaurantID, order.Status).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify after successful commit
	s.notification.NotifyRestaurant(order.RestaurantID, &order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
func (s *Server) handlePollOrder(w http.ResponseWriter, r *http.Request) {
	// Extract and validate restaurant ID
	restaurantID := strings.TrimPrefix(r.URL.Path, "/poll/orders/")
	restaurantIDInt, err := strconv.ParseInt(restaurantID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid restaurant ID", http.StatusBadRequest)
		return
	}

	// Subscribe to updates
	updates := s.notification.SubscribeRestaurant(restaurantIDInt)
	defer s.notification.UnsubscribeRestaurant(restaurantIDInt, updates)

	// Set headers for streaming response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case updatedOrder := <-updates:
		json.NewEncoder(w).Encode(updatedOrder)

	case <-r.Context().Done():
		// Client disconnected
		return

	case <-time.After(50 * time.Second):
		response := map[string]string{
			"status":        "timeout",
			"restaurant_id": restaurantID,
		}
		json.NewEncoder(w).Encode(response)
	}
}
