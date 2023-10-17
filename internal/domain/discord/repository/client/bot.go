package client

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka/internal/errors"
)

// JoinVoiceChannel to join voice channel.
func (c *Client) JoinVoiceChannel(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	// Already in voice channel.
	if c.voices[g.ID].isInVoiceChannel {
		return nil
	}

	// Looks for the user who call the command in voice channels.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			// Join voice channel.
			vc, err := c.session.ChannelVoiceJoin(g.ID, vs.ChannelID, false, false)
			if err != nil {
				return stack.Wrap(ctx, err)
			}

			c.voices[g.ID].Lock()
			c.voices[g.ID].voice = vc
			c.voices[g.ID].channelID = m.ChannelID
			c.voices[g.ID].isInVoiceChannel = true
			c.voices[g.ID].Unlock()

			break
		}
	}

	if !c.voices[g.ID].isInVoiceChannel {
		return stack.Wrap(ctx, errors.ErrNotInVC)
	}

	return nil
}

// LeaveVoiceChannel to leave voice channel.
func (c *Client) LeaveVoiceChannel(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	// Not in voice channel.
	if !c.voices[g.ID].isInVoiceChannel {
		return nil
	}

	// Leave voice channel.
	if err := c.voices[g.ID].voice.Disconnect(); err != nil {
		return stack.Wrap(ctx, err)
	}

	c.voices[g.ID].Lock()
	c.voices[g.ID].isInVoiceChannel = false
	c.voices[g.ID].Unlock()

	return nil
}
