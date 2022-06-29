package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/domain/template/entity"
	"github.com/rl404/naka/internal/errors"
)

// HandleQueue to handle queue.
func (s *service) HandleQueue(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	// Show queue list.
	if len(args) == 0 {
		queue := s.queue.GetList(ctx, g.ID)
		i := s.queue.GetIndex(ctx, g.ID)

		result := make([]entity.Video, len(queue))
		for i, q := range queue {
			result[i] = entity.Video{
				Title: q.Title,
			}
		}

		_, err := s.discord.SendMessage(ctx, m.ChannelID, s.template.GetQueue(i, result))
		return errors.Wrap(ctx, err)

	}

	// Search song.
	return errors.Wrap(ctx, s.searchSong(ctx, m, g, args, false))
}
