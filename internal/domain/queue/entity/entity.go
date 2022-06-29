package entity

import "time"

// Queue is entity for queue.
type Queue struct {
	Index int
	Songs []Song
}

// Song is entity for song.
type Song struct {
	Title        string
	URL          string
	ChannelTitle string
	ChannelURL   string
	Image        string
	Duration     time.Duration
	View         int
	Like         int
	Dislike      int
	SourceURL    string
}
