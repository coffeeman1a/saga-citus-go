package models

type SagaEvent struct {
	TxID    string  `json:"tx_id"`
	UserID  string  `json:"user_id"`
	Item    string  `json:"item"`
	Price   float64 `json:"price"`
	Step    string  `json:"step"`
	Status  string  `json:"status"`
	OrderID string  `json:"order_id,omitempty"`
}
