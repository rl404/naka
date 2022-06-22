package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/constant"
	"github.com/rl404/naka/internal/utils"
)

type template struct {
	prefix string
}

func newTemplate(prefix string) *template {
	return &template{
		prefix: prefix,
	}
}

func (t *template) clean(str string) string {
	return strings.ReplaceAll(str, "{{prefix}}", t.prefix)
}

func (t *template) getHelp() *discordgo.MessageEmbed {
	otherCmds := [][]string{
		{constant.MsgJoinCmd, constant.MsgJoinContent},
		{constant.MsgLeaveCmd, constant.MsgLeaveContent},
		{constant.MsgPauseCmd, constant.MsgPauseContent},
		{constant.MsgResumeCmd, constant.MsgResumeContent},
		{constant.MsgNextCmd, constant.MsgNextContent},
		{constant.MsgPrevCmd, constant.MsgPrevContent},
		{constant.MsgStopCmd, constant.MsgStopContent},
		{constant.MsgQueueCmd, constant.MsgQueueContent},
		{constant.MsgPurgeCmd, constant.MsgPurgeContent},
	}

	body := "```"
	for _, cmd := range otherCmds {
		body += utils.PadRight(t.clean(cmd[0]), len(t.prefix)+12, " ") + cmd[1] + "\n"
	}
	body += "```"

	return &discordgo.MessageEmbed{
		Title:       "Help",
		Description: constant.MsgHelpContent,
		Color:       constant.ColorGreyLight,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  constant.MsgPlayCmd,
				Value: t.clean(constant.MsgPlayContent),
			},
			{
				Name:  constant.MsgOtherCmd,
				Value: body,
			},
		},
	}
}

func (t *template) getQueue(queue []*audio) *discordgo.MessageEmbed {
	str := "```\n"
	for _, q := range queue {
		str += q.title + "\n"
	}
	str += "```"

	return &discordgo.MessageEmbed{
		Title:       "Queue",
		Description: str,
		Color:       constant.ColorBlue,
	}
}

func (t *template) playing(audio *audio) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       audio.title,
		Description: audio.duration.String(),
		Color:       constant.ColorBlue,
		URL:         audio.url,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "playing",
		},
	}
}
