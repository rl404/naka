package service

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/domain/template/entity"
	"github.com/rl404/naka/internal/errors"
)

// HandlePrompt to handle prompt response.
func (s *service) HandlePrompt(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, content string) error {
	search := s.prompt.GetSearch(ctx, m.Author.ID)
	if search != nil && len(search.IDs) > 0 {
		defer s.prompt.DeleteSearch(ctx, m.Author.ID)

		i, err := strconv.Atoi(content)
		if err != nil {
			if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidPrompt); err != nil {
				return errors.Wrap(ctx, err)
			}
			return errors.Wrap(ctx, err)
		}

		if i <= 0 || i > len(search.IDs) {
			if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidPrompt); err != nil {
				return errors.Wrap(ctx, err)
			}
			return errors.Wrap(ctx, errors.ErrInvalidPrompt)
		}

		videoID := search.IDs[i-1]

		return errors.Wrap(ctx, s.searchSong(ctx, m, g, []string{s.youtube.GenerateVideoURL(videoID)}, search.Play))
	}

	return nil
}
