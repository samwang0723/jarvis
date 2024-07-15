package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	config "github.com/samwang0723/jarvis/configs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.TimestampFieldName = "t"
	config.Load()
	cfg := config.GetCurrentConfig()
	logger := zerolog.New(os.Stdout).With().Str("app", "database-migration").Timestamp().Logger()

	ctx = logger.WithContext(ctx)

	migrationUp := true
	if config.IsMigrationDown() {
		migrationUp = false
	}

	runMigrate(ctx, migrationUp, cfg)
}
