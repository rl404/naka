package client

import (
	"context"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
)

// Client is discord client.
type Client struct {
	sync.Mutex

	session *discordgo.Session
	voices  map[string]*voice
}

// New to create new discord client.
func New(token string) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Client{
		session: session,
		voices:  make(map[string]*voice),
	}, nil
}

// Run to login and start discord bot.
func (c *Client) Run() error {
	return c.session.Open()
}

// Close to stop discord bot.
func (c *Client) Close() error {
	return c.session.Close()
}

// AddReadyHandler to add ready handler.
func (c *Client) AddReadyHandler(handler func(*discordgo.Session, *discordgo.Ready)) {
	c.session.AddHandler(handler)
}

// AddMessageHandler to add message handler.
func (c *Client) AddMessageHandler(handler func(*discordgo.Session, *discordgo.MessageCreate)) {
	c.session.AddHandler(handler)
}

// AddReactionHandler to add reaction handler.
func (c *Client) AddReactionHandler(handler func(*discordgo.Session, *discordgo.MessageReactionAdd)) {
	c.session.AddHandler(handler)
}

// GetGuildByChannelID to get guild by channel id.
func (c *Client) GetGuildByChannelID(ctx context.Context, cID string) (*discordgo.Guild, error) {
	ch, err := c.session.State.Channel(cID)
	if err != nil {
		return nil, stack.Wrap(ctx, err)
	}

	g, err := c.session.State.Guild(ch.GuildID)
	if err != nil {
		return nil, stack.Wrap(ctx, err)
	}

	return g, nil
}

// SendMessage to send message.
func (c *Client) SendMessage(ctx context.Context, cID, content string) (string, error) {
	m, err := c.session.ChannelMessageSend(cID, content)
	if err != nil {
		return "", stack.Wrap(ctx, err)
	}
	return m.ID, nil
}

// SendMessageEmbed to send embed message.
func (c *Client) SendMessageEmbed(ctx context.Context, cID string, content *discordgo.MessageEmbed) (string, error) {
	m, err := c.session.ChannelMessageSendEmbed(cID, content)
	if err != nil {
		return "", stack.Wrap(ctx, err)
	}
	return m.ID, nil
}

// EditMessage to edit message.
func (c *Client) EditMessage(ctx context.Context, cID, mID, content string) (string, error) {
	m, err := c.session.ChannelMessageEdit(cID, mID, content)
	if err != nil {
		return "", stack.Wrap(ctx, err)
	}
	return m.ID, nil
}

// EditMessageEmbed to edit embed message.
func (c *Client) EditMessageEmbed(ctx context.Context, cID, mID string, content *discordgo.MessageEmbed) (string, error) {
	m, err := c.session.ChannelMessageEditEmbed(cID, mID, content)
	if err != nil {
		return "", stack.Wrap(ctx, err)
	}
	return m.ID, nil
}
