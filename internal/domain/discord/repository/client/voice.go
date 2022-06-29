package client

import (
	"context"
	"io"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rl404/naka/internal/errors"
)

type voice struct {
	sync.Mutex

	voice            *discordgo.VoiceConnection
	channelID        string
	isInVoiceChannel bool
	isPlayerExist    bool
	isPlaying        bool
	isPaused         bool
	isStopped        bool
	disableAutoNext  bool

	encodeSession *dca.EncodeSession
	streamSession *dca.StreamingSession
}

// InitVoice to init discord voice.
func (c *Client) InitVoice(gID string) {
	c.Lock()
	defer c.Unlock()
	if c.voices[gID] == nil {
		c.voices[gID] = new(voice)
	}
}

// IsPlayerExist to get is player exist.
func (c *Client) IsPlayerExist(gID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.voices[gID].isPlayerExist
}

// SetPlayerExist to set player exist.
func (c *Client) SetPlayerExist(gID string, value bool) {
	c.Lock()
	defer c.Unlock()
	c.voices[gID].isPlayerExist = value
}

// GetChannelID to get channel id.
func (c *Client) GetChannelID(gID string) string {
	c.Lock()
	defer c.Unlock()
	return c.voices[gID].channelID
}

// GetStopped to get stopped.
func (c *Client) GetStopped(gID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.voices[gID].isStopped
}

// GetDisableAutoNext to get disable auto next.
func (c *Client) GetDisableAutoNext(gID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.voices[gID].disableAutoNext
}

// SetDisableAutoNext to set disable auto next.
func (c *Client) SetDisableAutoNext(gID string, value bool) {
	c.Lock()
	defer c.Unlock()
	c.voices[gID].disableAutoNext = value
}

// IsPlaying to get is playing.
func (c *Client) IsPlaying(gID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.voices[gID].isPlaying
}

// Pause to pause.
func (c *Client) Pause(gID string) {
	c.Lock()
	defer c.Unlock()

	if !c.voices[gID].isInVoiceChannel ||
		!c.voices[gID].isPlayerExist ||
		!c.voices[gID].isPlaying ||
		c.voices[gID].isPaused ||
		c.voices[gID].isStopped {
		return
	}

	c.voices[gID].isPaused = true
	c.voices[gID].streamSession.SetPaused(true)
}

// Resume to resume.
func (c *Client) Resume(gID string) {
	c.Lock()
	defer c.Unlock()

	if !c.voices[gID].isInVoiceChannel ||
		!c.voices[gID].isPlayerExist ||
		!c.voices[gID].isPlaying ||
		!c.voices[gID].isPaused ||
		c.voices[gID].isStopped {
		return
	}

	c.voices[gID].isPaused = false
	c.voices[gID].streamSession.SetPaused(false)
}

// SetPlaying to set playing.
func (c *Client) SetPlaying(ctx context.Context, gID string, value bool) error {
	c.Lock()
	defer c.Unlock()
	c.voices[gID].isPlaying = value
	return errors.Wrap(ctx, c.voices[gID].voice.Speaking(value))
}

// Stop to stop.
func (c *Client) Stop(gID string) {
	c.Lock()
	defer c.Unlock()

	if !c.voices[gID].isInVoiceChannel ||
		!c.voices[gID].isPlayerExist ||
		!c.voices[gID].isPlaying {
		return
	}

	c.voices[gID].disableAutoNext = true
	c.voices[gID].isStopped = true
	if c.voices[gID].encodeSession != nil {
		c.voices[gID].encodeSession.Cleanup()
	}
}

// Skip to skip.
func (c *Client) Skip(gID string) {
	c.Lock()
	defer c.Unlock()
	if c.voices[gID].encodeSession != nil {
		c.voices[gID].encodeSession.Cleanup()
	}
}

// Stream to stream audio.
func (c *Client) Stream(ctx context.Context, gID, path string) (err error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	c.voices[gID].encodeSession, err = dca.EncodeFile(path, options)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	defer c.voices[gID].encodeSession.Cleanup()

	done := make(chan error)

	c.voices[gID].streamSession = dca.NewStream(c.voices[gID].encodeSession, c.voices[gID].voice, done)

	err = <-done
	if err != nil && err != io.EOF {
		return errors.Wrap(ctx, err)
	}

	return nil
}
