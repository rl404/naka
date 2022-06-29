package repository

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Repository contains functions for discord domain.
type Repository interface {
	Run() error
	Close() error

	AddReadyHandler(func(*discordgo.Session, *discordgo.Ready))
	AddMessageHandler(func(*discordgo.Session, *discordgo.MessageCreate))
	AddReactionHandler(func(*discordgo.Session, *discordgo.MessageReactionAdd))

	SendMessage(ctx context.Context, channelID, content string) (string, error)
	SendMessageEmbed(ctx context.Context, channelID string, content *discordgo.MessageEmbed) (string, error)
	EditMessage(ctx context.Context, channelID, messageID, content string) (string, error)
	EditMessageEmbed(ctx context.Context, channelID, messageID string, content *discordgo.MessageEmbed) (string, error)

	InitVoice(guildID string)
	GetChannelID(guildID string) string

	IsPlayerExist(guildID string) bool
	SetPlayerExist(guildID string, value bool)
	GetStopped(guildID string) bool
	GetDisableAutoNext(guildID string) bool
	SetDisableAutoNext(guildID string, value bool)
	SetPlaying(ctx context.Context, guildID string, value bool) error
	Pause(guildID string)
	Resume(guildID string)
	Stop(guildID string)
	Skip(guildID string)

	Stream(ctx context.Context, guildID, path string) error

	JoinVoiceChannel(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	LeaveVoiceChannel(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error

	GetGuildByChannelID(ctx context.Context, channelID string) (*discordgo.Guild, error)
}
