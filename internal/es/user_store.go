package es

import (
	"context"
	"processor/internal/domain"
)

type UsersStore interface {
	IndexBulk(ctx context.Context, users []domain.User) error
	GetAll(ctx context.Context) ([]domain.User, error)
}
