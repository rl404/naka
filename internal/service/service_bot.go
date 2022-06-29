package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Run to run discord bot.
func (s *service) Run() error {
	return s.discord.Run()
}

// Stop to stop discord bot.
func (s *service) Stop() error {
	return s.discord.Close()
}

// GetGuildByChannelID to get guild by channel id.
func (s *service) GetGuildByChannelID(ctx context.Context, cID string) (*discordgo.Guild, error) {
	return s.discord.GetGuildByChannelID(ctx, cID)
}

// InitVoice to init discord voice.
func (s *service) InitVoice(guildID string) {
	s.discord.InitVoice(guildID)
}
