package handler

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/cache"
	"github.com/rl404/naka/internal/constant"
	"github.com/rl404/naka/internal/utils"
	"github.com/rl404/naka/internal/youtube"
)

type messageHandler struct {
	sync.Mutex

	cache    cache.Cacher
	prefix   string
	template *template
	youtube  *youtube.Client
	voices   map[string]*voice
}

// NewMessageHandler to create new discord message handler.
func NewMessageHandler(c cache.Cacher, prefix string, youtube *youtube.Client) func(*discordgo.Session, *discordgo.MessageCreate) {
	h := &messageHandler{
		cache:    c,
		prefix:   prefix,
		template: newTemplate(prefix),
		voices:   make(map[string]*voice),
		youtube:  youtube,
	}
	return h.handler()
}

func (h *messageHandler) handler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Get guild data.
		g, err := h.getGuildByChannelID(s, m.ChannelID)
		if err != nil {
			utils.Error(err.Error())
			return
		}

		if h.voices[g.ID] == nil {
			h.Lock()
			h.voices[g.ID] = new(voice)
			h.Unlock()
		}

		// Handle response for search song.
		if err := h.handleSearchResponse(s, m, g); err != nil {
			utils.Error(err.Error())
			return
		}

		// Command and prefix check.
		if h.prefixCheck(m.Content) {
			return
		}

		// Remove prefix.
		m.Content = h.cleanPrefix(m.Content)

		// Get arguments.
		r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
		args := r.FindAllString(m.Content, -1)

		switch args[0] {
		case "ping":
			err = h.handlePing(s, m)
		case "help", "h":
			err = h.handleHelp(s, m)
		case "join", "j":
			err = h.handleJoin(s, m, g)
		case "leave", "l":
			err = h.handleLeave(s, m, g)
		case "play", "p":
			err = h.handlePlay(s, m, g, args)
		case "pause":
			err = h.handlePause(s, m, g)
		case "resume":
			err = h.handleResume(s, m, g)
		case "skip", "next":
			err = h.handleNext(s, m, g)
		case "previous", "prev":
			err = h.handlePrevious(s, m, g)
		case "queue", "q":
			err = h.handleQueue(s, m, g, args)
		case "remove":
			err = h.handleRemove(s, m, g, args)
		case "stop":
			err = h.handleStop(s, m, g)
		case "purge":
			err = h.handlePurge(s, m, g)
		default:
			return
		}

		if err != nil {
			utils.Error(err.Error())
		}
	}
}

func (h *messageHandler) prefixCheck(cmd string) bool {
	return len(cmd) <= len(h.prefix) || cmd[:len(h.prefix)] != h.prefix
}

func (h *messageHandler) cleanPrefix(cmd string) string {
	return strings.TrimSpace(cmd[len(h.prefix):])
}

func (h *messageHandler) getGuildByChannelID(s *discordgo.Session, channelID string) (*discordgo.Guild, error) {
	c, err := s.State.Channel(channelID)
	if err != nil {
		return nil, err
	}
	return s.State.Guild(c.GuildID)
}

func (h *messageHandler) handlePing(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "pong")
	return err
}

func (h *messageHandler) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, h.template.getHelp())
	return err
}

func (h *messageHandler) handleJoin(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	return h.joinVoiceChannel(s, m, g)
}

func (h *messageHandler) joinVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	// Already joined voice channel.
	if h.voices[g.ID].isInVoiceChannel {
		return nil
	}

	// Looks for the user who call the command in voice channels.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			// Join voice channel.
			vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, false)
			if err != nil {
				return err
			}

			h.voices[g.ID].Lock()
			h.voices[g.ID].voice = vc
			h.voices[g.ID].session = s
			h.voices[g.ID].guildID = g.ID
			h.voices[g.ID].channelID = m.ChannelID
			h.voices[g.ID].isInVoiceChannel = true
			h.voices[g.ID].template = h.template
			h.voices[g.ID].Unlock()

			break
		}
	}

	// User not in voice channel.
	if h.voices[g.ID].voice == nil || !h.voices[g.ID].isInVoiceChannel {
		if _, err := s.ChannelMessageSend(m.ChannelID, constant.MsgNotInVC); err != nil {
			return err
		}
	}

	return nil
}

func (h *messageHandler) handleLeave(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	return h.leaveVoiceChannel(s, m, g)
}

func (h *messageHandler) leaveVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	if !h.voices[g.ID].isInVoiceChannel {
		return nil
	}

	// Stop first.
	h.voices[g.ID].stop()

	// Leave voice channel.
	if err := h.voices[g.ID].voice.Disconnect(); err != nil {
		return err
	}

	h.voices[g.ID].Lock()
	h.voices[g.ID].isInVoiceChannel = false
	h.voices[g.ID].Unlock()

	return nil
}

func (h *messageHandler) handlePlay(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	// Just play the queue.
	if len(args) == 1 {
		// Join voice channel.
		if err := h.joinVoiceChannel(s, m, g); err != nil {
			return err
		}

		go h.voices[g.ID].play()
		return nil
	}

	// Search.
	return h.searchMusic(s, m, g, args[1:], true)
}

func (h *messageHandler) searchMusic(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild, args []string, play bool) error {
	// Add to queue if given youtube url.
	if h.youtube.IsYoutubeLink(args[0]) {
		id, err := h.youtube.GetVideoIDFromURL(args[0])
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, constant.MsgInvalidYoutube); err != nil {
				return err
			}
			return err
		}

		path, err := h.youtube.GetPath(id)
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, constant.MsgInvalidYoutube); err != nil {
				return err
			}
			return err
		}

		video, err := h.youtube.GetVideo(id)
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, constant.MsgInvalidYoutube); err != nil {
				return err
			}
			return err
		}

		h.voices[g.ID].addQueue(audio{
			title: video.Title,
			path:  path,
			url:   h.youtube.GenerateLink(id),
			image: video.Image,
		})

		if play {
			if err := h.joinVoiceChannel(s, m, g); err != nil {
				return err
			}

			go h.voices[g.ID].play()
		}

		return nil
	}

	videos, err := h.youtube.Search(strings.Join(args, " "), 10)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, err.Error()); err != nil {
			return err
		}
		return err
	}

	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, h.template.search(videos)); err != nil {
		return err
	}

	if len(videos) == 0 {
		return nil
	}

	key := utils.GetKey("search", m.Author.ID)
	if err := h.cache.Set(key, videos); err != nil {
		return err
	}

	return nil
}

func (h *messageHandler) handleSearchResponse(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	var videos []youtube.Video
	key := utils.GetKey("search", m.Author.ID)
	if err := h.cache.Get(key, &videos); err != nil {
		return nil
	}

	if err := h.cache.Delete(key); err != nil {
		return err
	}

	// Select search result no.
	i, err := strconv.Atoi(m.Content)
	if err != nil {
		return nil
	}

	if i <= 0 || i > len(videos) {
		return nil
	}

	id := videos[i-1].ID

	path, err := h.youtube.GetPath(id)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, constant.MsgInvalidYoutube); err != nil {
			return err
		}
		return err
	}

	video, err := h.youtube.GetVideo(id)
	if err != nil {
		if _, err := s.ChannelMessageSend(m.ChannelID, constant.MsgInvalidYoutube); err != nil {
			return err
		}
		return err
	}

	h.voices[g.ID].addQueue(audio{
		title: video.Title,
		path:  path,
		url:   h.youtube.GenerateLink(id),
		image: video.Image,
	})

	if err := h.joinVoiceChannel(s, m, g); err != nil {
		return err
	}

	go h.voices[g.ID].play()
	return nil
}

func (h *messageHandler) handlePause(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	if !h.voices[g.ID].isInVoiceChannel {
		return nil
	}

	h.voices[g.ID].pause()

	currentAudio := h.voices[g.ID].getCurrentAudio()

	if _, err := h.voices[g.ID].session.ChannelMessageSendEmbed(h.voices[g.ID].channelID, h.voices[g.ID].template.paused(currentAudio)); err != nil {
		utils.Error(err.Error())
		return err
	}

	return nil
}

func (h *messageHandler) handleResume(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	if !h.voices[g.ID].isInVoiceChannel {
		return nil
	}

	h.voices[g.ID].resume()

	currentAudio := h.voices[g.ID].getCurrentAudio()

	if _, err := h.voices[g.ID].session.ChannelMessageSendEmbed(h.voices[g.ID].channelID, h.voices[g.ID].template.playing(currentAudio)); err != nil {
		utils.Error(err.Error())
		return err
	}

	return nil
}

func (h *messageHandler) handleNext(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	h.voices[g.ID].next()
	return nil
}

func (h *messageHandler) handlePrevious(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	h.voices[g.ID].previous()
	return nil
}

func (h *messageHandler) handleQueue(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	if len(args) == 1 {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, h.template.getQueue(h.voices[g.ID].queue))
		return err
	}

	return h.searchMusic(s, m, g, args[1:], false)
}

func (h *messageHandler) handleStop(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	if !h.voices[g.ID].isInVoiceChannel {
		return nil
	}

	h.voices[g.ID].stop()

	if _, err := h.voices[g.ID].session.ChannelMessageSend(h.voices[g.ID].channelID, "stopped"); err != nil {
		return err
	}

	return nil
}

func (h *messageHandler) handlePurge(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild) error {
	h.voices[g.ID].purgeQueue()

	if _, err := h.voices[g.ID].session.ChannelMessageSend(h.voices[g.ID].channelID, "queue purged"); err != nil {
		return err
	}

	return nil
}

func (h *messageHandler) handleRemove(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild, args []string) error {
	if len(args) == 1 {
		if _, err := h.voices[g.ID].session.ChannelMessageSend(m.ChannelID, constant.MsgInvalid); err != nil {
			return err
		}
	}

	for _, arg := range args[1:] {
		i, err := strconv.Atoi(arg)
		if err != nil {
			if _, err := h.voices[g.ID].session.ChannelMessageSend(m.ChannelID, constant.MsgInvalid); err != nil {
				return err
			}
			return nil
		}

		h.voices[g.ID].deleteQueue(i - 1)
	}

	return nil
}
