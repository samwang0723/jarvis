package main

import (
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/server"
	"github.com/samwang0723/jarvis/internal/helper"
)

func main() {
	config.Load()
	cfg := config.GetCurrentConfig()
	logger := zerolog.New(os.Stdout).With().Str("app", cfg.Server.Name).Timestamp().Logger()

	// Set the global log level based on the environment variable
	switch strings.ToLower(cfg.Log.Level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // Default to Info level
	}

	// manually set time zone, docker image may not have preset timezone
	var err error
	time.Local, err = time.LoadLocation(helper.TimeZone)
	if err != nil {
		logger.Error().Msgf("error loading location '%s': %v\n", helper.TimeZone, err)
	}

	server.Serve(cfg, &logger)
}
