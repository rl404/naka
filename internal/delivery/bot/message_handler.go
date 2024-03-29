package bot

import (
	"context"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
)

func (b *Bot) messageHandler(nrApp *newrelic.Application) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		ctx := stack.Init(context.Background())
		defer b.log(ctx)

		// Ignore all messages created by the bot itself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Get guild.
		g, err := b.service.GetGuildByChannelID(ctx, m.ChannelID)
		if err != nil {
			stack.Wrap(ctx, err)
			return
		}

		// Init voice.
		b.service.InitVoice(g.ID)

		// Handle prompt.
		if err := b.service.HandlePrompt(ctx, m, g, m.Content); err != nil {
			stack.Wrap(ctx, err)
			return
		}

		// Command and prefix check.
		if b.prefixCheck(m.Content) {
			return
		}

		// Remove prefix.
		m.Content = b.cleanPrefix(m.Content)

		// Get arguments.
		r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
		args := r.FindAllString(m.Content, -1)

		tx := nrApp.StartTransaction("Command " + args[0])
		defer tx.End()

		ctx = newrelic.NewContext(ctx, tx)

		switch args[0] {
		case "ping":
			stack.Wrap(ctx, b.service.HandlePing(ctx, m))
		case "help", "h":
			stack.Wrap(ctx, b.service.HandleHelp(ctx, m))
		case "play", "p":
			stack.Wrap(ctx, b.service.HandlePlay(ctx, m, g, args[1:]))
		case "join", "j":
			stack.Wrap(ctx, b.service.HandleJoin(ctx, m, g))
		case "leave", "l":
			stack.Wrap(ctx, b.service.HandleLeave(ctx, m, g))
		case "pause":
			stack.Wrap(ctx, b.service.HandlePause(ctx, m, g))
		case "resume":
			stack.Wrap(ctx, b.service.HandleResume(ctx, m, g))
		case "stop", "s":
			stack.Wrap(ctx, b.service.HandleStop(ctx, m, g))
		case "next":
			stack.Wrap(ctx, b.service.HandleNext(ctx, m, g))
		case "previous", "prev":
			stack.Wrap(ctx, b.service.HandlePrev(ctx, m, g))
		case "skip", "jump":
			stack.Wrap(ctx, b.service.HandleJump(ctx, m, g, args[1:]))
		case "queue", "q":
			stack.Wrap(ctx, b.service.HandleQueue(ctx, m, g, args[1:]))
		case "remove":
			stack.Wrap(ctx, b.service.HandleRemove(ctx, m, g, args[1:]))
		case "purge":
			stack.Wrap(ctx, b.service.HandlePurge(ctx, m, g))
		}
	}
}

func (b *Bot) prefixCheck(cmd string) bool {
	return len(cmd) <= len(b.prefix) || cmd[:len(b.prefix)] != b.prefix
}

func (b *Bot) cleanPrefix(cmd string) string {
	return strings.TrimSpace(cmd[len(b.prefix):])
}
