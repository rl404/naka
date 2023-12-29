package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka/internal/domain/template/entity"
)

// HandleQueue to handle queue.
func (s *service) HandleQueue(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	// Show queue list.
	if len(args) == 0 {
		queue := s.queue.GetList(ctx, g.ID)
		currentIndex := s.queue.GetIndex(ctx, g.ID)

		result := make([]entity.Video, len(queue))

		resultIndex := 0
		for i := currentIndex; i < len(queue); i++ {
			result[resultIndex] = entity.Video{
				Title: queue[i].Title,
			}

			resultIndex++
		}

		_, err := s.discord.SendMessage(ctx, m.ChannelID, s.template.GetQueue(currentIndex, result))
		return stack.Wrap(ctx, err)

	}

	// Search song.
	return stack.Wrap(ctx, s.searchSong(ctx, m, g, args, false))
}
