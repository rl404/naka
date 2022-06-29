package client

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/rl404/naka/internal/domain/youtube/entity"
	"github.com/rl404/naka/internal/errors"
	"github.com/rl404/naka/internal/utils"
	"google.golang.org/api/googleapi/transport"
	_youtube "google.golang.org/api/youtube/v3"
)

// Client is youtube client.
type Client struct {
	client  *youtube.Client
	service *_youtube.Service
}

// New to create new youtube client.
func New(key string) (*Client, error) {
	service, err := _youtube.New(&http.Client{
		Transport: &transport.APIKey{Key: key},
		Timeout:   5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		client: &youtube.Client{
			HTTPClient: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
		service: service,
	}, nil
}

// GenerateVideoURL to generate youtube video url.
func (c *Client) GenerateVideoURL(id string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
}

// GenerateChannelURL to generate youtube channel url.
func (c *Client) GenerateChannelURL(id string) string {
	return fmt.Sprintf("https://www.youtube.com/channel/%s", id)
}

// IsURLValid to check if youtube url valid.
func (c *Client) IsURLValid(url_ string) bool {
	if _, err := url.ParseRequestURI(url_); err != nil {
		return false
	}
	id, err := c.GetIDFromURL(context.Background(), url_)
	return id != "" && err == nil
}

// GetIDFromURL to get video id from url.
func (c *Client) GetIDFromURL(ctx context.Context, url string) (string, error) {
	id, err := youtube.ExtractVideoID(url)
	if err != nil {
		return "", errors.Wrap(ctx, errors.ErrInvalidYoutubeURL, err)
	}
	return id, nil
}

// GetSourceURLByID to get video source url.
func (c *Client) GetSourceURLByID(ctx context.Context, id string) (string, error) {
	video, err := c.client.GetVideoContext(ctx, id)
	if err != nil {
		return "", errors.Wrap(ctx, err)
	}

	format := video.Formats.WithAudioChannels()

	url, err := c.client.GetStreamURLContext(ctx, video, &format[0])
	if err != nil {
		return "", errors.Wrap(ctx, err)
	}

	return url, nil
}

// GetVideos to search youtube video.
func (c *Client) GetVideos(ctx context.Context, query string, limit int64) ([]entity.Video, error) {
	response, err := c.service.Search.
		List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(limit).
		Order("viewCount").
		Context(ctx).
		Do()
	if err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	res := make([]entity.Video, len(response.Items))
	for i, data := range response.Items {
		res[i] = entity.Video{
			ID:    data.Id.VideoId,
			Title: strings.TrimSpace(html.UnescapeString(data.Snippet.Title)),
		}
	}

	return res, nil
}

// GetVideo to get video detail.
func (c *Client) GetVideo(ctx context.Context, id string) (*entity.Video, error) {
	response, err := c.service.Videos.
		List([]string{"id", "snippet", "contentDetails", "statistics"}).
		Id(id).
		Context(ctx).
		Do()
	if err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	if len(response.Items) == 0 {
		return nil, errors.Wrap(ctx, errors.ErrInvalidYoutubeID)
	}

	return &entity.Video{
		ID:           response.Items[0].Id,
		Title:        strings.TrimSpace(html.UnescapeString(response.Items[0].Snippet.Title)),
		ChannelID:    response.Items[0].Snippet.ChannelId,
		ChannelTitle: strings.TrimSpace(html.UnescapeString(response.Items[0].Snippet.ChannelTitle)),
		Image:        response.Items[0].Snippet.Thumbnails.Default.Url,
		Duration:     utils.ParseDuration(response.Items[0].ContentDetails.Duration),
		View:         int(response.Items[0].Statistics.ViewCount),
		Like:         int(response.Items[0].Statistics.LikeCount),
		Dislike:      int(response.Items[0].Statistics.DislikeCount),
	}, nil
}
