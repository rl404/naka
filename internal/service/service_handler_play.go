package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/domain/template/entity"
	"github.com/rl404/naka/internal/errors"
)

// HandlePlay to handle play.
func (s *service) HandlePlay(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	// Just play the queue.
	if len(args) == 0 {
		return s.play(ctx, m, g)
	}

	// Search song.
	return errors.Wrap(ctx, s.searchSong(ctx, m, g, args, true))
}

func (s *service) play(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	// Player already exist.
	if s.discord.IsPlayerExist(g.ID) {
		return nil
	}

	// Check queue.
	if s.queue.IsEmpty(ctx, g.ID) {
		if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.EmptyQueue); err != nil {
			return errors.Wrap(ctx, err)
		}
		return nil
	}

	// Join channel.
	if err := s.HandleJoin(ctx, m, g); err != nil {
		return errors.Wrap(ctx, err)
	}

	// Loop the queue.
	go func() error {
		s.discord.SetPlayerExist(g.ID, true)
		defer s.discord.SetPlayerExist(g.ID, false)

		for {
			// Get queue.
			queue := s.queue.GetList(ctx, g.ID)
			i := s.queue.GetIndex(ctx, g.ID)

			if len(queue) == 0 || i >= len(queue) {
				if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.EndQueue); err != nil {
					return errors.Wrap(ctx, err)
				}
				return nil
			}

			song := queue[i]

			if _, err := s.discord.SendMessageEmbed(ctx, m.ChannelID, s.template.GetPlaying(entity.Video{
				Title:        song.Title,
				URL:          song.URL,
				ChannelTitle: song.ChannelTitle,
				ChannelURL:   song.ChannelURL,
				Image:        song.Image,
				Duration:     song.Duration,
				View:         song.View,
				Like:         song.Like,
				Dislike:      song.Dislike,
				QueueI:       i + 1,
				QueueCnt:     len(queue),
			})); err != nil {
				return errors.Wrap(ctx, err)
			}

			if err := s.discord.SetPlaying(ctx, g.ID, true); err != nil {
				return errors.Wrap(ctx, err)
			}

			// Start stream.
			if err := s.discord.Stream(ctx, g.ID, song.SourceURL); err != nil {
				return errors.Wrap(ctx, err)
			}

			// Go next queue.
			if !s.discord.GetDisableAutoNext(g.ID) {
				s.queue.SetIndex(ctx, g.ID, i+1)
			}
			s.discord.SetDisableAutoNext(g.ID, false)

			if err := s.discord.SetPlaying(ctx, g.ID, false); err != nil {
				return errors.Wrap(ctx, err)
			}

			if s.discord.GetStopped(g.ID) {
				return nil
			}
		}
	}()

	return nil
}
