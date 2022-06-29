package repository

import (
	"context"

	"github.com/rl404/naka/internal/domain/prompt/entity"
)

// Repository contains functions for action domain.
type Repository interface {
	SetSearch(ctx context.Context, userID string, data entity.Search) error
	GetSearch(ctx context.Context, userID string) *entity.Search
	DeleteSearch(ctx context.Context, userID string) error
}
