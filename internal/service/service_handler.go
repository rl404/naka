package service

import (
	"context"
	_errors "errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/domain/template/entity"
	"github.com/rl404/naka/internal/errors"
)

// RegisterReadyHandler to register discord ready handler.
func (s *service) RegisterReadyHandler(fn func(*discordgo.Session, *discordgo.Ready)) {
	s.discord.AddReadyHandler(fn)
}

// RegisterMessageHandler to register discord message handler.
func (s *service) RegisterMessageHandler(fn func(*discordgo.Session, *discordgo.MessageCreate)) {
	s.discord.AddMessageHandler(fn)
}

// HandlePing to handle ping.
func (s *service) HandlePing(ctx context.Context, m *discordgo.MessageCreate) error {
	_, err := s.discord.SendMessage(ctx, m.ChannelID, "pong")
	return errors.Wrap(ctx, err)
}

// HandleHelp to handle help.
func (s *service) HandleHelp(ctx context.Context, m *discordgo.MessageCreate) error {
	_, err := s.discord.SendMessageEmbed(ctx, m.ChannelID, s.template.GetHelp())
	return errors.Wrap(ctx, err)
}

// HandleJoin to handle join.
func (s *service) HandleJoin(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	err := s.discord.JoinVoiceChannel(ctx, m, g)
	if _errors.Is(err, errors.ErrNotInVC) {
		if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.NotInVC); err != nil {
			return errors.Wrap(ctx, err)
		}
	}
	return errors.Wrap(ctx, err)
}

// HandleLeave to handle leave.
func (s *service) HandleLeave(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	if err := s.HandleStop(ctx, m, g); err != nil {
		return errors.Wrap(ctx, err)
	}
	return errors.Wrap(ctx, s.discord.LeaveVoiceChannel(ctx, m, g))
}

// HandlePause to handle pause.
func (s *service) HandlePause(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	s.discord.Pause(g.ID)
	_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.Paused)
	return errors.Wrap(ctx, err)
}

// HandleResume to handle resume.
func (s *service) HandleResume(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	s.discord.Resume(g.ID)
	_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.Resumed)
	return errors.Wrap(ctx, err)
}

// HandleStop to handle stop.
func (s *service) HandleStop(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	s.discord.Stop(g.ID)
	_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.Stopped)
	return errors.Wrap(ctx, err)
}

// HandleNext to handle next.
func (s *service) HandleNext(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	i := s.queue.GetIndex(ctx, g.ID)
	s.queue.SetIndex(ctx, g.ID, i+1)
	s.discord.SetDisableAutoNext(g.ID, true)
	s.discord.Skip(g.ID)

	_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.Next)
	return errors.Wrap(ctx, err)
}

// HandlePrev to handle previous.
func (s *service) HandlePrev(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	i := s.queue.GetIndex(ctx, g.ID)
	s.queue.SetIndex(ctx, g.ID, i-1)
	s.discord.SetDisableAutoNext(g.ID, true)
	s.discord.Skip(g.ID)

	_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.Previous)
	return errors.Wrap(ctx, err)
}

// HandleJump to handle jump.
func (s *service) HandleJump(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	if len(args) == 0 {
		_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidQueue)
		return errors.Wrap(ctx, err)
	}

	ii, err := strconv.Atoi(args[0])
	if err != nil {
		_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidQueue)
		return errors.Wrap(ctx, err)
	}

	s.queue.SetIndex(ctx, g.ID, ii-1)
	s.discord.SetDisableAutoNext(g.ID, true)
	s.discord.Skip(g.ID)

	_, err = s.discord.SendMessage(ctx, m.ChannelID, s.template.GetJumped(ii))
	return errors.Wrap(ctx, err)
}

// HandleRemove to handle remove.
func (s *service) HandleRemove(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	if len(args) == 0 {
		_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidQueue)
		return errors.Wrap(ctx, err)
	}

	is := make([]int, len(args))
	for i, arg := range args {
		j, err := strconv.Atoi(arg)
		if err != nil {
			if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidQueue); err != nil {
				return errors.Wrap(ctx, err)
			}
			return errors.Wrap(ctx, err)
		}

		is[i] = j - 1
	}

	if err := s.queue.Remove(ctx, g.ID, is...); err != nil {
		return errors.Wrap(ctx, err)
	}

	_, err := s.discord.SendMessage(ctx, m.ChannelID, s.template.GetRemoved(args))
	return errors.Wrap(ctx, err)
}

// HandlePurge to handle purge.
func (s *service) HandlePurge(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	if err := s.queue.Purge(ctx, g.ID); err != nil {
		return errors.Wrap(ctx, err)
	}

	_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.Purged)
	return errors.Wrap(ctx, err)
}
