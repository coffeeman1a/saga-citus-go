package models

import "time"

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Item      string    `json:"item"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type NewOrderRequest struct {
	UserID string  `json:"user_id"`
	Item   string  `json:"item"`
	Price  float64 `json:"price"`
}

type StatusUpdateRquest struct {
	Status string `json:"status"`
}
