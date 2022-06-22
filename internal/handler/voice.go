package handler

import (
	"io"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rl404/naka/internal/constant"
	"github.com/rl404/naka/internal/utils"
)

type voice struct {
	sync.Mutex
	template *template

	session          *discordgo.Session
	voice            *discordgo.VoiceConnection
	guildID          string
	channelID        string
	isInVoiceChannel bool
	isPlayerExist    bool

	queueI        int
	queue         []*audio
	encodeSession *dca.EncodeSession
	streamSession *dca.StreamingSession
	isPlaying     bool
	isPaused      bool
	isStopped     bool
	autoNext      bool
}

type audio struct {
	title    string
	duration time.Duration
	url      string
	path     string
}

func (v *voice) getQueue() []*audio {
	v.Lock()
	defer v.Unlock()
	return v.queue
}

func (v *voice) addQueue(audio audio) {
	v.Lock()
	defer v.Unlock()
	v.queue = append(v.queue, &audio)
}

func (v *voice) getCurrentAudio() *audio {
	v.Lock()
	defer v.Unlock()
	return v.queue[v.queueI]
}

func (v *voice) nextQueue() {
	v.Lock()
	defer v.Unlock()
	if v.queueI <= len(v.queue)-1 {
		v.queueI++
	}
}

func (v *voice) previousQueue() {
	v.Lock()
	defer v.Unlock()
	if v.queueI > 0 {
		v.queueI--
	}
}

func (v *voice) deleteQueue() {
	v.Lock()
	defer v.Unlock()
	v.queueI, v.queue = 0, nil
}

func (v *voice) next() {
	v.autoNext = false

	if v.queueI >= len(v.queue)-1 {
		v.stop()
		return
	}

	if v.encodeSession != nil {
		v.encodeSession.Cleanup()
	}

	v.nextQueue()
}

func (v *voice) previous() {
	v.autoNext = false

	if v.queueI <= 0 {
		v.stop()
		return
	}

	if v.encodeSession != nil {
		v.encodeSession.Cleanup()
	}

	v.previousQueue()
}

func (v *voice) stop() {
	v.autoNext = false

	if !v.isPlaying {
		return
	}

	v.isStopped = true
	if v.encodeSession != nil {
		v.encodeSession.Cleanup()
	}
}

func (v *voice) pause() {
	if !v.isPlaying || v.isPaused {
		return
	}

	v.Lock()
	defer v.Unlock()

	v.isPaused, v.isPlaying = true, false
	v.streamSession.SetPaused(true)
	if err := v.voice.Speaking(false); err != nil {
		utils.Error(err.Error())
	}
}

func (v *voice) resume() {
	if v.isPlaying || !v.isPaused {
		return
	}

	v.Lock()
	defer v.Unlock()

	v.isPaused, v.isPlaying = false, true
	v.streamSession.SetPaused(false)
	if err := v.voice.Speaking(true); err != nil {
		utils.Error(err.Error())
	}
}

func (v *voice) play() error {
	// Empty queue.
	if len(v.queue) == 0 {
		_, err := v.session.ChannelMessageSend(v.channelID, constant.MsgEmptyQueue)
		return err
	}

	// Still playing.
	if v.isPlayerExist {
		return nil
	}

	v.isPlayerExist, v.autoNext, v.isPaused, v.isStopped = true, true, false, false
	defer func() { v.isPlayerExist = false }()

	for {
		if len(v.queue) == 0 {
			if _, err := v.session.ChannelMessageSend(v.channelID, constant.MsgEndQueue); err != nil {
				utils.Error(err.Error())
				return err
			}
			return nil
		}

		currentAudio := v.getCurrentAudio()

		if _, err := v.session.ChannelMessageSendEmbed(v.channelID, v.template.playing(currentAudio)); err != nil {
			utils.Error(err.Error())
			return err
		}

		if err := v.voice.Speaking(true); err != nil {
			utils.Error(err.Error())
			return err
		}
		v.isPlaying = true

		// Start stream.
		if err := v.stream(currentAudio.path); err != nil {
			utils.Error(err.Error())
			return err
		}

		// Go next queue.
		if v.autoNext {
			v.nextQueue()
		}
		v.autoNext = true

		v.isPlaying = false
		if err := v.voice.Speaking(false); err != nil {
			utils.Error(err.Error())
			return err
		}

		if v.isStopped {
			return nil
		}
	}
}

func (v *voice) stream(path string) (err error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	v.encodeSession, err = dca.EncodeFile(path, options)
	if err != nil {
		return err
	}

	defer v.encodeSession.Cleanup()

	done := make(chan error)

	v.streamSession = dca.NewStream(v.encodeSession, v.voice, done)

	err = <-done
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
