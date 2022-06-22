package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/naka/internal/discord"
	"github.com/rl404/naka/internal/handler"
	"github.com/rl404/naka/internal/utils"
	"github.com/rl404/naka/internal/youtube"
)

func bot() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Init cache.
	c, err := cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
	if err != nil {
		return err
	}
	utils.Info("cache initialized")
	defer c.Close()

	// Init youtube.
	youtube := youtube.New(cfg.Youtube.Key)

	// Init discord.
	discord, err := discord.New(cfg.Discord.Token)
	if err != nil {
		return err
	}
	utils.Info("repository discord initialized")

	// Init & add handler.
	discord.AddReadyHandler(handler.NewReadyHandler(cfg.Discord.Prefix))
	discord.AddMessageHandler(handler.NewMessageHandler(c, cfg.Discord.Prefix, youtube))

	// Run bot.
	if err := discord.Run(); err != nil {
		return err
	}
	utils.Info("naka is running...")
	defer discord.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	return nil
}
