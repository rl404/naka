package repository

import (
	"context"

	"github.com/rl404/naka/internal/domain/queue/entity"
)

// Repository contains functions for queue domain.
type Repository interface {
	GetList(ctx context.Context, guildID string) []entity.Song
	IsEmpty(ctx context.Context, guildID string) bool
	Add(ctx context.Context, guildID string, data entity.Song) error
	Remove(ctx context.Context, guildID string, i ...int) error
	Purge(ctx context.Context, guildID string) error

	GetIndex(ctx context.Context, guildID string) int
	SetIndex(ctx context.Context, guildID string, i int) error
}
