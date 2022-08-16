package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/naka/internal/domain/queue/entity"
	"github.com/rl404/naka/internal/errors"
	"github.com/rl404/naka/internal/utils"
)

// Cache contains functions for queue cache.
type Cache struct {
	cacher cache.Cacher
}

// New to create new queue cache.
func New(cacher cache.Cacher) *Cache {
	return &Cache{
		cacher: cacher,
	}
}

func (c *Cache) getQueue(ctx context.Context, gID string) (data *entity.Queue) {
	key := utils.GetKey("queue", gID)
	c.cacher.Get(ctx, key, &data)
	return data
}

func (c *Cache) setQueue(ctx context.Context, gID string, data entity.Queue) error {
	key := utils.GetKey("queue", gID)
	return errors.Wrap(ctx, c.cacher.Set(ctx, key, data))
}

// GetList to get list.
func (c *Cache) GetList(ctx context.Context, gID string) []entity.Song {
	q := c.getQueue(ctx, gID)
	if q == nil {
		return nil
	}
	return q.Songs
}

// IsEmpty to check if empty.
func (c *Cache) IsEmpty(ctx context.Context, gID string) bool {
	q := c.getQueue(ctx, gID)
	return q == nil || len(q.Songs) == 0
}

// Add to add song to queue.
func (c *Cache) Add(ctx context.Context, gID string, data entity.Song) error {
	q := c.getQueue(ctx, gID)
	if q == nil {
		q = &entity.Queue{}
	}
	q.Songs = append(q.Songs, data)
	return errors.Wrap(ctx, c.setQueue(ctx, gID, *q))
}

// Remove to remove song from queue.
func (c *Cache) Remove(ctx context.Context, gID string, is ...int) error {
	q := c.getQueue(ctx, gID)
	if q == nil {
		q = &entity.Queue{}
	}

	newQ := entity.Queue{Index: q.Index}

	for i := range q.Songs {
		for _, j := range is {
			if i == j {
				q.Songs[i] = entity.Song{Title: ""}
			}
		}
	}

	for _, s := range q.Songs {
		if s.Title != "" {
			newQ.Songs = append(newQ.Songs, s)
		}
	}

	return errors.Wrap(ctx, c.setQueue(ctx, gID, newQ))
}

// Purge to delete queue.
func (c *Cache) Purge(ctx context.Context, gID string) error {
	key := utils.GetKey("queue", gID)
	return errors.Wrap(ctx, c.cacher.Delete(ctx, key))
}

// GetIndex to get index.
func (c *Cache) GetIndex(ctx context.Context, gID string) int {
	q := c.getQueue(ctx, gID)
	if q == nil {
		return 0
	}
	if q.Index > len(q.Songs) {
		q.Index = len(q.Songs)
	}
	if q.Index < 0 {
		q.Index = 0
	}
	return q.Index
}

// SetIndex to set index.
func (c *Cache) SetIndex(ctx context.Context, gID string, i int) error {
	q := c.getQueue(ctx, gID)
	if q == nil {
		q = &entity.Queue{}
	}
	q.Index = i
	if i > len(q.Songs) {
		q.Index = len(q.Songs)
	}
	if i < 0 {
		q.Index = 0
	}
	return errors.Wrap(ctx, c.setQueue(ctx, gID, *q))
}
