package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // pgx driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/database"
)

type MigrateOption struct {
	DB    *config.Config
	Up    bool
	Reset bool
}

func runMigrate(ctx context.Context, migrationUp bool, dbConfig *config.Config) {
	logger := zerolog.Ctx(ctx).With().Logger()

	source, err := iofs.New(database.MigrationFiles, "migrations")
	if err != nil {
		logger.Fatal().Err(err).Msg("initial iofs failed")
	}

	dbURL := dbURL(dbConfig)

	dbm, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("initial migrate instance failed")
	}

	defer dbm.Close()

	dbm.Log = &migrateLogger{&logger}

	action := "up"
	if !migrationUp {
		action = "down"
	}

	logger = logger.With().Str("action", action).Logger()

	logger.Info().Msg("database migration start!")

	var dbmErr error

	if migrationUp {
		dbmErr = dbm.Up()
	} else {
		dbmErr = dbm.Down()
	}

	if dbmErr != nil {
		if !errors.Is(dbmErr, migrate.ErrNoChange) {
			logger.Error().Err(dbmErr).Msg("database migration failed")

			return
		}

		logger.Info().Msg("database migration no change!")
	}

	logger.Info().Msg("database migration SUCCESS!")
}

type migrateLogger struct {
	logger *zerolog.Logger
}

func (ml *migrateLogger) Printf(format string, v ...any) {
	ml.logger.Info().Msg(fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...))
}

func (ml *migrateLogger) Verbose() bool {
	return true
}

func dbURL(cfg *config.Config) string {
	return fmt.Sprintf(
		"pgx5://%s:%s@%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		net.JoinHostPort(cfg.Database.Host, cfg.Database.Port),
		cfg.Database.Database,
	)
}
