package repository

import (
	"context"
	"errors"

	"github.com/coffeeman1a/saga-citus-go/users-service/internal/models"

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

func (r *PGXRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.Pool.QueryRow(ctx,
		`select id, email 
	from users where id = $1`, id).Scan(&user.ID, &user.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *PGXRepository) NewUser(ctx context.Context, email string) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.Pool.Exec(ctx,
		`insert into users (id, email) values ($1, $2)`,
		id, email)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
