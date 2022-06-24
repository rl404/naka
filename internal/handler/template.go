package handler

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/constant"
	"github.com/rl404/naka/internal/utils"
	"github.com/rl404/naka/internal/youtube"
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
	body := "```\n"
	for i, q := range queue {
		body += utils.PadLeft(strconv.Itoa(i+1), 2, " ") + " " + q.title + "\n"
	}
	body += "```"

	return &discordgo.MessageEmbed{
		Title:       "Queue",
		Description: body,
		Color:       constant.ColorBlue,
	}
}

func (t *template) playing(audio *audio) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: audio.title,
		Color: constant.ColorBlue,
		URL:   audio.url,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: audio.image,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: "playing",
		},
	}
}

func (t *template) search(videos []youtube.Video) *discordgo.MessageEmbed {
	var body string
	if len(videos) == 0 {
		body = "empty..."
	} else {
		body = "```\n"
		for i, v := range videos {
			body += utils.PadLeft(strconv.Itoa(i+1), 2, " ") + " " + v.Title + "\n"
		}
		body += "```"
	}

	return &discordgo.MessageEmbed{
		Title:       "Search Results",
		Description: body,
		Color:       constant.ColorBlue,
	}
}
