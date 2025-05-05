package repository

import (
	"context"

	"github.com/coffeeman1a/saga-citus-go/orders-service/internal/models"
	"github.com/google/uuid"
)

type OrderRepository interface {
	NewOrder(ctx context.Context, user_id string, item string, price float64) (uuid.UUID, error)
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status string) (bool, error)
}
