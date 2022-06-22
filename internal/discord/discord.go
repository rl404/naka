package discord

import "github.com/bwmarrin/discordgo"

type discord struct {
	client *discordgo.Session
}

// New to create new discord client.
func New(token string) (*discord, error) {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &discord{
		client: client,
	}, nil
}

// Run to login and start discord bot.
func (d *discord) Run() error {
	return d.client.Open()
}

// Close to stop discord bot.
func (d *discord) Close() error {
	return d.client.Close()
}

// AddReadyHandler to add ready handler.
func (d *discord) AddReadyHandler(handler func(*discordgo.Session, *discordgo.Ready)) {
	d.client.AddHandler(handler)
}

// AddMessageHandler to add message handler.
func (d *discord) AddMessageHandler(handler func(*discordgo.Session, *discordgo.MessageCreate)) {
	d.client.AddHandler(handler)
}

// AddReactionHandler to add reaction handler.
func (d *discord) AddReactionHandler(handler func(*discordgo.Session, *discordgo.MessageReactionAdd)) {
	d.client.AddHandler(handler)
}
