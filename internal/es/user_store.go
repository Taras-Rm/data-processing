package es

import (
	"context"
	"processor/internal/domain"
)

type UserStore interface {
	IndexBulk(ctx context.Context, users []domain.User) error
	Search(ctx context.Context) ([]domain.User, error)
}
