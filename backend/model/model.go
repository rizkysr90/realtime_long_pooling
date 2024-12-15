package model

type Order struct {
	ID           int64  `json:"id"`
	RestaurantID int64  `json:"restaurant_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

type Notification struct {
	ID           int64  `json:"id"`
	RestaurantID int64  `json:"restaurant_id"`
	OrderID      int64  `json:"order_id"`
	Type         string `json:"type"`
	ReadStatus   bool   `json:"read_status"`
	CreatedAt    string `json:"created_at"`
}
