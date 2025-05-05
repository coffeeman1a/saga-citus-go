package repository

import (
	"context"

	"github.com/coffeeman1a/saga-citus-go/users-service/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	NewUser(ctx context.Context, email string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}
