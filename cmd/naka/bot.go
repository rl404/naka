package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/cache"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrCache "github.com/rl404/fairy/monitoring/newrelic/cache"
	_bot "github.com/rl404/naka/internal/delivery/bot"
	discordRepository "github.com/rl404/naka/internal/domain/discord/repository"
	discordClient "github.com/rl404/naka/internal/domain/discord/repository/client"
	promptRepository "github.com/rl404/naka/internal/domain/prompt/repository"
	promptCache "github.com/rl404/naka/internal/domain/prompt/repository/cache"
	queueRepository "github.com/rl404/naka/internal/domain/queue/repository"
	queueCache "github.com/rl404/naka/internal/domain/queue/repository/cache"
	templateRepository "github.com/rl404/naka/internal/domain/template/repository"
	templateClient "github.com/rl404/naka/internal/domain/template/repository/client"
	youtubeRepository "github.com/rl404/naka/internal/domain/youtube/repository"
	youtubeClient "github.com/rl404/naka/internal/domain/youtube/repository/client"
	"github.com/rl404/naka/internal/service"
	"github.com/rl404/naka/internal/utils"
)

func bot() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Init newrelic.
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Newrelic.Name),
		newrelic.ConfigLicense(cfg.Newrelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		utils.Error(err.Error())
	} else {
		defer nrApp.Shutdown(10 * time.Second)
		utils.AddLog(_nr.NewFromNewrelicApp(nrApp, _nr.ErrorLevel))
		utils.Info("newrelic initialized")
	}

	// Init cache.
	c, err := cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
	if err != nil {
		return err
	}
	c = nrCache.New(cfg.Cache.Dialect, c)
	utils.Info("cache initialized")
	defer c.Close()

	// Init discord.
	var discord discordRepository.Repository
	discord, err = discordClient.New(cfg.Discord.Token)
	if err != nil {
		return err
	}
	utils.Info("discord initialized")

	// Init youtube.
	var youtube youtubeRepository.Repository
	youtube, err = youtubeClient.New(cfg.Youtube.Key)
	if err != nil {
		return err
	}
	utils.Info("youtube initialized")

	// Init template.
	var template templateRepository.Repository = templateClient.New(cfg.Discord.Prefix)
	utils.Info("template initialized")

	// Init queue.
	var queue queueRepository.Repository = queueCache.New(c)
	utils.Info("queue initialized")

	// Init prompt.
	var prompt promptRepository.Repository = promptCache.New(c)
	utils.Info("prompt initialized")

	// Init service.
	service := service.New(discord, youtube, template, queue, prompt)
	utils.Info("service initialized")

	// Init bot.
	bot := _bot.New(service, cfg.Discord.Prefix)
	bot.RegisterHandler(nrApp)
	utils.Info("bot initialized")

	// Run bot.
	if err := bot.Run(); err != nil {
		return err
	}
	utils.Info("naka is running...")
	defer bot.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	return nil
}
