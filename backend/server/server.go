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

	switch {
	case r.Method == "POST" && r.URL.Path == "/orders":
		s.handlerCreateOrder(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/poll/orders/"): // Added leading slash
		s.handlePollOrder(w, r)
	// case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/orders/"):
	// 	// s.handleGetOrder(w, r)
	// case r.Method == "PATCH" && strings.HasSuffix(r.URL.Path, "/status"):
	// 	// s.handleUpdateOrderStatus(w, r)
	// case r.Method == "GET" && r.URL.Path == "/notifications":
	// 	// s.handleGetNotifications(w, r)
	// case r.Method == "PATCH" && strings.HasSuffix(r.URL.Path, "/read"):
	// 	// s.handleMarkNotificationRead(w, r)
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
	query := `
        INSERT INTO orders (restaurant_id, status)
        VALUES ($1, $2)
        RETURNING id, created_at`

	err := s.db.QueryRow(query, order.RestaurantID, "new").Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.notification.NotifyRestaurant(order.RestaurantID, &order)
	json.NewEncoder(w).Encode(order)
}
func (s *Server) handlePollOrder(w http.ResponseWriter, r *http.Request) {
	restaurantID := strings.TrimPrefix(r.URL.Path, "/poll/orders/")
	restaurantIDInt, _ := strconv.Atoi(restaurantID)
	// Subscribe to updates
	updates := s.notification.SubscribeRestaurant(int64(restaurantIDInt))
	defer s.notification.UnscubscribeRestaurant(int64(restaurantIDInt), updates)

	select {
	case updatedJob := <-updates:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 OK

		// Marshal and write the JSON response
		json.NewEncoder(w).Encode(updatedJob)
	case <-r.Context().Done():
		// Client disconnected
		return
	case <-time.After(50 * time.Second):
		response := map[string]interface{}{
			"restaurant_id": restaurantIDInt,
			"status":        "timeout",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
