package youtube

import (
	"strings"

	"github.com/kkdai/youtube/v2"
)

// Client is youtube client.
type Client struct {
	key    string
	client *youtube.Client
}

// New to create new youtube client.
func New(key string) *Client {
	return &Client{
		key:    key,
		client: &youtube.Client{},
	}
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
func (c *Client) GetURL(id string) (string, error) {
	video, err := c.client.GetVideo(id)
	if err != nil {
		return "", err
	}

	format := video.Formats.WithAudioChannels()

	return c.client.GetStreamURL(video, &format[0])
}
