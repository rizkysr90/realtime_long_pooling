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
	n.Mu.Lock()
	defer n.Mu.Unlock()

	if subscribers, exists := n.OrderListeners[restaurantID]; exists {
		for _, subscriber := range subscribers {
			select {
			case subscriber <- data:
				// Successfully sent update
			default:
				// Channel is blocked, skip this listener
			}
		}
	}
}
func (n *Notification) SubscribeRestaurant(restaurantID int64) chan *model.Order {
	// Thread safe to handle race condition
	n.Mu.Lock()
	defer n.Mu.Unlock()

	// create new subscriber
	newSubscriber := make(chan *model.Order, 1)
	n.OrderListeners[restaurantID] = append(n.OrderListeners[restaurantID], newSubscriber)
	return newSubscriber
}
func (n *Notification) UnscubscribeRestaurant(restaurantID int64, ch chan *model.Order) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	if subscribers, exists := n.OrderListeners[restaurantID]; exists {
		for i, subs := range subscribers {
			if subs == ch {
				// Remove this channel from Listeners
				n.OrderListeners[restaurantID] = append(
					n.OrderListeners[restaurantID][:i], n.OrderListeners[restaurantID][i+1:]...)
				close(ch)
				break
			}
		}
	}
}
