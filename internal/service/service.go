package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	discordRepository "github.com/rl404/naka/internal/domain/discord/repository"
	promptRepository "github.com/rl404/naka/internal/domain/prompt/repository"
	queueRepository "github.com/rl404/naka/internal/domain/queue/repository"
	templateRepository "github.com/rl404/naka/internal/domain/template/repository"
	youtubeRepository "github.com/rl404/naka/internal/domain/youtube/repository"
)

// Service contains functions for service.
type Service interface {
	Run() error
	Stop() error

	RegisterReadyHandler(func(*discordgo.Session, *discordgo.Ready))
	RegisterMessageHandler(func(*discordgo.Session, *discordgo.MessageCreate))

	GetGuildByChannelID(ctx context.Context, channelID string) (*discordgo.Guild, error)

	InitVoice(guildID string)

	HandlePing(ctx context.Context, m *discordgo.MessageCreate) error
	HandleHelp(ctx context.Context, m *discordgo.MessageCreate) error
	HandlePrompt(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, content string) error
	HandlePlay(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error
	HandleJoin(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandleLeave(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandlePause(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandleResume(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandleStop(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandleNext(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandlePrev(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
	HandleJump(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error
	HandleQueue(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error
	HandleRemove(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error
	HandlePurge(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error
}

type service struct {
	discord  discordRepository.Repository
	youtube  youtubeRepository.Repository
	template templateRepository.Repository
	queue    queueRepository.Repository
	prompt   promptRepository.Repository
}

// New to create new service.
func New(
	discord discordRepository.Repository,
	youtube youtubeRepository.Repository,
	template templateRepository.Repository,
	queue queueRepository.Repository,
	prompt promptRepository.Repository,
) Service {
	return &service{
		discord:  discord,
		youtube:  youtube,
		template: template,
		queue:    queue,
		prompt:   prompt,
	}
}
