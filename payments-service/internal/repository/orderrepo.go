package repository

import (
	"context"

	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/models"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	NewPayment(ctx context.Context, order_id string, amount float64) (uuid.UUID, error)
	GetPaymentByID(ctx context.Context, id string) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id string, status string) (bool, error)
}
