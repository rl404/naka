package errors

import "errors"

// List of error.
var (
	ErrInvalidCacheTime  = errors.New("invalid cache time")
	ErrInvalidYoutubeURL = errors.New("invalid youtube url")
	ErrInvalidYoutubeID  = errors.New("invalid youtube video id")
	ErrNotInVC           = errors.New("not in voice channel")
	ErrInvalidPrompt     = errors.New("invalid prompt response")
)
