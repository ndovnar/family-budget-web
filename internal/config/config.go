package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/ndovnar/family-budget-api/internal/api"
	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/store/mongo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel zerolog.Level `default:"info" desc:"Level for generated logs"`
	Auth     auth.Config
	API      api.Config
	Mongo    mongo.Config
}

func Load() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("CONFIG", &cfg)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(cfg.LogLevel)

	return cfg, err
}
