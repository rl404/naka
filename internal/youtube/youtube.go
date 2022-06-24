package youtube

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/kkdai/youtube/v2"
	"google.golang.org/api/googleapi/transport"
	_youtube "google.golang.org/api/youtube/v3"
)

// Client is youtube client.
type Client struct {
	key     string
	host    string
	client  *youtube.Client
	service *_youtube.Service
}

// Video is youtube video model.
type Video struct {
	ID    string
	Title string
	Image string
}

// New to create new youtube client.
func New(key string) (*Client, error) {
	service, err := _youtube.New(&http.Client{
		Transport: &transport.APIKey{Key: key},
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		key:     key,
		client:  &youtube.Client{},
		service: service,
	}, nil
}

// GenerateLink to generate youtube link.
func (c *Client) GenerateLink(id string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
}

// IsYoutubeLink to check if url is youtube.
func (c *Client) IsYoutubeLink(url string) bool {
	return strings.Contains(url, "https://www.youtube.com")
}

// GetVideoIDFromURL to get video id from url.
func (c *Client) GetVideoIDFromURL(url string) (string, error) {
	return youtube.ExtractVideoID(url)
}

// GetPath to get video path.
func (c *Client) GetPath(id string) (string, error) {
	video, err := c.client.GetVideo(id)
	if err != nil {
		return "", err
	}

	format := video.Formats.WithAudioChannels()

	return c.client.GetStreamURL(video, &format[0])
}

// Search to search youtube video.
func (c *Client) Search(query string, limit int64) ([]Video, error) {
	response, err := c.service.Search.
		List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(limit).
		Do()
	if err != nil {
		return nil, err
	}

	res := make([]Video, len(response.Items))
	for i, data := range response.Items {
		res[i] = Video{
			ID:    data.Id.VideoId,
			Title: html.UnescapeString(data.Snippet.Title),
			Image: data.Snippet.Thumbnails.Default.Url,
		}
	}

	return res, nil
}

// GetVideo to get video detail.
func (c *Client) GetVideo(id string) (*Video, error) {
	response, err := c.service.Videos.
		List([]string{"id", "snippet"}).
		Id(id).
		Do()
	if err != nil {
		return nil, err
	}

	return &Video{
		ID:    response.Items[0].Id,
		Title: response.Items[0].Snippet.Title,
		Image: response.Items[0].Snippet.Thumbnails.Default.Url,
	}, nil
}
