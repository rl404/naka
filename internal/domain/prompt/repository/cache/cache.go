package cache

import (
	"context"
	"time"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/naka/internal/domain/prompt/entity"
	"github.com/rl404/naka/internal/errors"
	"github.com/rl404/naka/internal/utils"
)

// Cache contains functions for prompt cache.
type Cache struct {
	cacher cache.Cacher
}

// New to create new prompt cache.
func New(cacher cache.Cacher) *Cache {
	return &Cache{
		cacher: cacher,
	}
}

// SetSearch to set search prompt.
func (c *Cache) SetSearch(ctx context.Context, userID string, data entity.Search) error {
	key := utils.GetKey("prompt", "search", userID)
	if err := c.cacher.Set(ctx, key, data, time.Minute); err != nil {
		return errors.Wrap(ctx, err)
	}
	return nil
}

// GetSearch to get search prompt.
func (c *Cache) GetSearch(ctx context.Context, userID string) (data *entity.Search) {
	key := utils.GetKey("prompt", "search", userID)
	c.cacher.Get(ctx, key, &data)
	return data
}

// DeleteSearch to delete search prompt.
func (c *Cache) DeleteSearch(ctx context.Context, userID string) error {
	key := utils.GetKey("prompt", "search", userID)
	return errors.Wrap(ctx, c.cacher.Delete(ctx, key))
}
