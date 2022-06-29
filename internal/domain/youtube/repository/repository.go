package repository

import (
	"context"

	"github.com/rl404/naka/internal/domain/youtube/entity"
)

// Repository contains functions for youtube domain.
type Repository interface {
	GenerateVideoURL(id string) string
	GenerateChannelURL(id string) string
	IsURLValid(url string) bool
	GetIDFromURL(ctx context.Context, url string) (string, error)
	GetSourceURLByID(ctx context.Context, id string) (string, error)
	GetVideos(ctx context.Context, query string, limit int64) ([]entity.Video, error)
	GetVideo(ctx context.Context, id string) (*entity.Video, error)
}
