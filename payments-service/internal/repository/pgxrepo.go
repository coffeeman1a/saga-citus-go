package repository

import (
	"context"
	"errors"

	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/config"
	"github.com/coffeeman1a/saga-citus-go/payments-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXRepository struct {
	Pool *pgxpool.Pool
}

func NewPGXRepository(p *pgxpool.Pool) *PGXRepository {
	return &PGXRepository{Pool: p}
}

func (r *PGXRepository) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	var payment models.Payment
	err := r.Pool.QueryRow(ctx,
		`select id, order_id, status, created_at
	from payments where id = $1`, id).Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

func (r *PGXRepository) NewPayment(ctx context.Context, order_id string, amount float64) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.Pool.Exec(ctx,
		`insert into payments (id, order_id, amount, status) values ($1, $2, $3, $4)`,
		id, order_id, amount, config.StatusReserved)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PGXRepository) UpdatePaymentStatus(ctx context.Context, id string, status string) (bool, error) {
	if status == "" {
		return false, errors.New("empty status string")
	}

	cmdTag, err := r.Pool.Exec(ctx,
		`update payments
		set status = $1
	where id = $2`, status, id)

	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, errors.New("payment not found")
	}

	return true, nil
}
