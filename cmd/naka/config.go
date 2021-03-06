package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/log"
	"github.com/rl404/naka/internal/errors"
	"github.com/rl404/naka/internal/utils"
)

type config struct {
	Discord discordConfig `envconfig:"DISCORD"`
	Cache   cacheConfig   `envconfig:"CACHE"`
	Youtube youtubeConfig `envconfig:"YOUTUBE"`
	Log     logConfig     `envconfig:"LOG"`
}

type discordConfig struct {
	Token  string `envconfig:"TOKEN" required:"true"`
	Prefix string `envconfig:"PREFIX" required:"true" default:"="`
}

type cacheConfig struct {
	Dialect  string        `envconfig:"DIALECT" required:"true" default:"inmemory"`
	Address  string        `envconfig:"ADDRESS"`
	Password string        `envconfig:"PASSWORD"`
	Time     time.Duration `envconfig:"TIME" required:"true" default:"24h"`
}

type youtubeConfig struct {
	Key string `envconfig:"KEY" required:"true"`
}

type logConfig struct {
	Type  log.LogType  `envconfig:"TYPE" default:"2"`
	Level log.LogLevel `envconfig:"LEVEL" default:"-1"`
	JSON  bool         `envconfig:"JSON" default:"false"`
	Color bool         `envconfig:"COLOR" default:"true"`
}

const envPath = "../../.env"
const envPrefix = "NAKA"

var cacheType = map[string]cache.CacheType{
	"nocache":  cache.NoCache,
	"redis":    cache.Redis,
	"inmemory": cache.InMemory,
	"memcache": cache.Memcache,
}

func getConfig() (*config, error) {
	var cfg config

	// Load .env file.
	_ = godotenv.Load(envPath)

	// Convert env to struct.
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}

	if cfg.Cache.Time <= 0 {
		return nil, errors.ErrInvalidCacheTime
	}

	// Init global log.
	if err := utils.InitLog(cfg.Log.Type, cfg.Log.Level, cfg.Log.JSON, cfg.Log.Color); err != nil {
		return nil, err
	}

	return &cfg, nil
}
