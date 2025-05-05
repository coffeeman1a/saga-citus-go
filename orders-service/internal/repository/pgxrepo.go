package repository

import (
	"context"
	"errors"

	"github.com/coffeeman1a/saga-citus-go/orders-service/internal/config"
	"github.com/coffeeman1a/saga-citus-go/orders-service/internal/models"

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

func (r *PGXRepository) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.Pool.QueryRow(ctx,
		`select id, user_id, item, price, status, created_at
	from orders where id = $1`, id).Scan(&order.ID, &order.UserID, &order.Item, &order.Price, &order.Status, &order.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (r *PGXRepository) NewOrder(ctx context.Context, user_id string, item string, price float64) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.Pool.Exec(ctx,
		`insert into orders (id, user_id, item, price, status) values ($1, $2, $3, $4, $5)`,
		id, user_id, item, price, config.StatusPending)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PGXRepository) UpdateOrderStatus(ctx context.Context, id string, status string) (bool, error) {
	if status == "" {
		return false, errors.New("empty status string")
	}

	cmdTag, err := r.Pool.Exec(ctx,
		`update orders
		set status = $1
	where id = $2`, status, id)

	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, errors.New("order not found")
	}

	return true, nil
}
