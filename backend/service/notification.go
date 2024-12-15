package service

import (
	"realtime-dashboard-food-delivery/model"
	"sync"
)

type Notification struct {
	Mu             sync.RWMutex
	OrderListeners map[int64][]chan *model.Order
}

func NewNotification() *Notification {
	return &Notification{
		Mu:             sync.RWMutex{},
		OrderListeners: make(map[int64][]chan *model.Order),
	}
}
func (n *Notification) NotifyRestaurant(restaurantID int64, data *model.Order) {
	// 1. Quick read with RLock
	n.Mu.RLock()
	var activeSubscribers []chan *model.Order
	if subscribers, exists := n.OrderListeners[restaurantID]; exists {
		activeSubscribers = make([]chan *model.Order, len(subscribers))
		copy(activeSubscribers, subscribers)
	}
	n.Mu.RUnlock()

	// 2. Send notifications without holding any lock
	for _, subscriber := range activeSubscribers {
		select {
		case subscriber <- data:
		default:
		}
	}
	// Send notifications outside the lock
	for _, subscriber := range activeSubscribers {
		select {
		case subscriber <- data:
			// Successfully sent update
		default:
			// Channel is blocked, skip this listener
			// Consider cleaning up blocked channels
		}
	}
}
func (n *Notification) SubscribeRestaurant(restaurantID int64) chan *model.Order {
	// Thread safe to handle race condition
	n.Mu.Lock()
	defer n.Mu.Unlock()
	// Create buffered channel to prevent blocking
	updates := make(chan *model.Order, 1)

	n.OrderListeners[restaurantID] = append(n.OrderListeners[restaurantID], updates)
	return updates
}
func (n *Notification) UnsubscribeRestaurant(restaurantID int64, ch chan *model.Order) {
	n.Mu.Lock()
	defer n.Mu.Unlock()

	subscribers := n.OrderListeners[restaurantID]
	for i, subscriber := range subscribers {
		if subscriber == ch {
			// Remove the subscriber
			n.OrderListeners[restaurantID] = append(subscribers[:i], subscribers[i+1:]...)
			close(ch)
			break
		}
	}

	// Clean up if no more subscribers
	if len(n.OrderListeners[restaurantID]) == 0 {
		delete(n.OrderListeners, restaurantID)
	}
}
