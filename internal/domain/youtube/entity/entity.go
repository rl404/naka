package entity

import "time"

// Video is youtube video entity.
type Video struct {
	ID           string
	Title        string
	ChannelID    string
	ChannelTitle string
	Image        string
	Duration     time.Duration
	View         int
	Like         int
	Dislike      int
}
