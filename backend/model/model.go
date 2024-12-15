package model

type Order struct {
	ID           int64  `json:"id"`
	RestaurantID int64  `json:"restaurant_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}
