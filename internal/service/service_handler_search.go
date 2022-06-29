package service

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	promptEntity "github.com/rl404/naka/internal/domain/prompt/entity"
	queueEntity "github.com/rl404/naka/internal/domain/queue/entity"
	"github.com/rl404/naka/internal/domain/template/entity"
	"github.com/rl404/naka/internal/errors"
)

func (s *service) searchSong(ctx context.Context, m *discordgo.MessageCreate, g *discordgo.Guild, args []string, play bool) error {
	if len(args) == 0 {
		_, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidSearchQuery)
		return errors.Wrap(ctx, err)
	}

	// If using youtube url.
	if s.youtube.IsURLValid(args[0]) {
		videoID, err := s.youtube.GetIDFromURL(ctx, args[0])
		if err != nil {
			if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidYoutubeURL); err != nil {
				return errors.Wrap(ctx, err)
			}
			return errors.Wrap(ctx, err)
		}

		sourceURL, err := s.youtube.GetSourceURLByID(ctx, videoID)
		if err != nil {
			if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidYoutubeURL); err != nil {
				return errors.Wrap(ctx, err)
			}
			return errors.Wrap(ctx, err)
		}

		video, err := s.youtube.GetVideo(ctx, videoID)
		if err != nil {
			if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidYoutubeURL); err != nil {
				return errors.Wrap(ctx, err)
			}
			return errors.Wrap(ctx, err)
		}

		if err := s.queue.Add(ctx, g.ID, queueEntity.Song{
			Title:        video.Title,
			URL:          s.youtube.GenerateVideoURL(videoID),
			ChannelTitle: video.ChannelTitle,
			ChannelURL:   s.youtube.GenerateChannelURL(video.ChannelID),
			Image:        video.Image,
			Duration:     video.Duration,
			View:         video.View,
			Like:         video.Like,
			Dislike:      video.Dislike,
			SourceURL:    sourceURL,
		}); err != nil {
			return errors.Wrap(ctx, err)
		}

		if _, err := s.discord.SendMessageEmbed(ctx, m.ChannelID, s.template.GetAddQueue(entity.Video{
			Title:        video.Title,
			URL:          s.youtube.GenerateVideoURL(videoID),
			ChannelTitle: video.ChannelTitle,
			ChannelURL:   s.youtube.GenerateChannelURL(video.ChannelID),
			Image:        video.Image,
			Duration:     video.Duration,
			View:         video.View,
			Like:         video.Like,
			Dislike:      video.Dislike,
			QueueI:       len(s.queue.GetList(ctx, g.ID)),
		})); err != nil {
			return errors.Wrap(ctx, err)
		}

		if play {
			return errors.Wrap(ctx, s.play(ctx, m, g))
		}

		return nil
	}

	// Search song in youtube.
	videos, err := s.youtube.GetVideos(ctx, strings.Join(args, " "), 10)
	if err != nil {
		if _, err := s.discord.SendMessage(ctx, m.ChannelID, entity.InvalidSearchQuery); err != nil {
			return errors.Wrap(ctx, err)
		}
		return errors.Wrap(ctx, err)
	}

	result := make([]entity.Video, len(videos))
	ids := make([]string, len(videos))
	for i, v := range videos {
		result[i] = entity.Video{Title: v.Title}
		ids[i] = v.ID
	}

	mID, err := s.discord.SendMessage(ctx, m.ChannelID, s.template.GetSearch(result))
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	if len(videos) == 0 {
		return nil
	}

	// Prepare prompt.
	if err := s.prompt.SetSearch(ctx, m.Author.ID, promptEntity.Search{
		Play: play,
		IDs:  ids,
	}); err != nil {
		return errors.Wrap(ctx, err)
	}

	go func() {
		// Remove prompt when expired.
		time.Sleep(time.Minute)
		s.discord.EditMessage(ctx, m.ChannelID, mID, entity.PromptExpired)
	}()

	return nil
}
