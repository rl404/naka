package repository

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka/internal/domain/template/entity"
)

// Repository contains functions for template.
type Repository interface {
	GetHelp() *discordgo.MessageEmbed
	GetPlaying(data entity.Video) *discordgo.MessageEmbed
	GetAddQueue(data entity.Video) *discordgo.MessageEmbed
	GetSearch(data []entity.Video) string
	GetQueue(i int, data []entity.Video) string
	GetJumped(i int) string
	GetRemoved(i []string) string
}
