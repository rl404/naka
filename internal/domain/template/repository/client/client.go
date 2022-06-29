package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/domain/template/entity"
	"github.com/rl404/naka/internal/utils"
)

// Client is template client.
type Client struct {
	prefix string
}

// New to create new template client.
func New(prefix string) *Client {
	return &Client{
		prefix: prefix,
	}
}

func (c *Client) clean(str string) string {
	return strings.ReplaceAll(str, "{{prefix}}", c.prefix)
}

// GetHelp to get help template.
func (c *Client) GetHelp() *discordgo.MessageEmbed {
	cmds := [][]string{
		{entity.JoinCmd, entity.JoinDesc},
		{entity.LeaveCmd, entity.LeaveDesc},
		{entity.PauseCmd, entity.PauseDesc},
		{entity.ResumeCmd, entity.ResumeDesc},
		{entity.NextCmd, entity.NextDesc},
		{entity.PrevCmd, entity.PrevDesc},
		{entity.StopCmd, entity.StopDesc},
		{entity.QueueCmd, entity.QueueDesc},
		{entity.RemoveCmd, entity.RemoveDesc},
		{entity.PurgeCmd, entity.PurgeDesc},
	}

	body := "```"
	for _, cmd := range cmds {
		body += fmt.Sprintf("%s%s\n",
			utils.PadRight(c.clean(cmd[0]), len(c.prefix)+12, " "),
			cmd[1],
		)
	}
	body += "```"

	return &discordgo.MessageEmbed{
		Title:       "Help",
		Description: entity.HelpDesc,
		Color:       entity.ColorGreyLight,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Play Song",
				Value: c.clean(entity.PlayDesc),
			},
			{
				Name:  "Other Commands",
				Value: body,
			},
		},
	}
}

// GetPlaying to get playing template.
func (c *Client) GetPlaying(data entity.Video) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: data.Title,
		Color: entity.ColorBlue,
		URL:   data.URL,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: data.Image,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: data.ChannelTitle,
			URL:  data.ChannelURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("playing queue %d/%d", data.QueueI, data.QueueCnt),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Duration",
				Value:  data.Duration.String(),
				Inline: true,
			},
			{
				Name:   "View",
				Value:  utils.Thousands(data.View),
				Inline: true,
			},
		},
	}
}

// GetAddQueue to get add queue template.
func (c *Client) GetAddQueue(data entity.Video) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: data.Title,
		Color: entity.ColorBlue,
		URL:   data.URL,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: data.Image,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: data.ChannelTitle,
			URL:  data.ChannelURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("queued at %d", data.QueueI),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Duration",
				Value:  data.Duration.String(),
				Inline: true,
			},
			{
				Name:   "View",
				Value:  utils.Thousands(data.View),
				Inline: true,
			},
		},
	}
}

// GetSearch to get search template.
func (c *Client) GetSearch(data []entity.Video) string {
	body := "**Search Results**\n"

	if len(data) == 0 {
		body += "not found..."
	} else {
		body += "```"
		for i, v := range data {
			body += fmt.Sprintf("%s | %s\n",
				utils.PadLeft(strconv.Itoa(i+1), 2, " "),
				v.Title,
			)
		}
		body += "```"
		body += "*Choose the song by typing the number.*"
	}

	return body
}

// GetQueue to get queue template.
func (c *Client) GetQueue(index int, data []entity.Video) string {
	body := "**Queue**\n"

	if len(data) == 0 {
		body = entity.EmptyQueue
	} else {
		body += "```"
		for i, v := range data {
			no := strconv.Itoa(i + 1)
			if i == index {
				no = "->"
			}

			body += fmt.Sprintf("%s | %s\n",
				utils.PadLeft(no, 2, " "),
				v.Title,
			)
		}
		body += "```"
	}

	return body
}

// GetJumped to get jumped template.
func (c *Client) GetJumped(i int) string {
	return fmt.Sprintf("Jumped to %d.", i)
}

// GetRemoved to get removed template.
func (c *Client) GetRemoved(i []string) string {
	return fmt.Sprintf("Removed %s from queue.", strings.Join(i, ", "))
}
