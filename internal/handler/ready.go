package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type readyHandler struct {
	prefix string
}

// NewReadyHandler to create new discord ready handler.
func NewReadyHandler(prefix string) func(*discordgo.Session, *discordgo.Ready) {
	h := &readyHandler{prefix: prefix}
	return h.handler()
}

func (h *readyHandler) handler() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		s.UpdateListeningStatus(fmt.Sprintf("%shelp for command list", h.prefix))
	}
}
