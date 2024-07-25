package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	config "github.com/samwang0723/jarvis/configs"

	database "github.com/samwang0723/jarvis/database"
)

const (
	migrationUp   = "up"
	migrationDown = "down"
)

//nolint:gochecknoglobals // only allowed global vars - filled at build time - do not change
var (
	CommitTime = "dev"
	CommitHash = "dev"
)

func main() {
	run()
}

func run() {
	env := config.GetCurrentEnv()

	// logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if env != "local" {
		os.Setenv("APP_ENV", env)
	}

	action := strings.ToLower(flag.Arg(0))
	if action == "" {
		action = migrationUp

		logger.Info().Msg("database migration arg not present, use default: " + migrationUp)
	}

	type ConfigDatabase struct {
		Port     string `env:"PSQL_PORT"   env-default:"DB_PORT"`
		Host     string `env:"PSQL_HOST"   env-default:"DB_HOST"`
		Name     string `env:"PSQL_DBNAME" env-default:"APP_NAME_UND_main"`
		User     string `env:"PSQL_USER"   env-default:"APP_NAME_UND_app"`
		Password string `env:"PSQL_PASS"   env-default:"DB_PASSWORD"`
	}

	var cfg ConfigDatabase

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		logger.Error().Err(err).Msg("env var load failed")
	}

	dsn := fmt.Sprintf("pgx://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		net.JoinHostPort(cfg.Host, cfg.Port),
		cfg.Name,
	)

	// migrate
	source, err := iofs.New(database.MigrationFiles, "migrations")
	if err != nil {
		logger.Fatal().Err(err).Msg("initial iofs failed")
	}

	dbm, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("initial migrate instance failed")
	}

	defer dbm.Close()

	dbm.Log = &migrateLogger{&logger}

	logger.Info().
		Str("commitTime", CommitTime).
		Str("commitHash", CommitHash).
		Msg("database migration " + action + " start!")

	var dbmErr error

	if action == migrationUp {
		dbmErr = dbm.Up()
	} else if action == migrationDown {
		dbmErr = dbm.Down()
	}

	if dbmErr != nil {
		if errors.Is(dbmErr, migrate.ErrNoChange) {
			logger.Info().Msg("database migration " + action + " no change!")
			logger.Info().
				Msg("database migration " + action + " SUCCESS!")
			// the log `SUCCESS` used for circleci detection
		} else {
			logger.Error().Err(dbmErr).Msg("database migration " + action + " failed")
		}
	} else {
		logger.Info().Msg("database migration " + action + " SUCCESS!") // the log `SUCCESS` used for circleci detection
	}
}

type migrateLogger struct {
	logger *zerolog.Logger
}

func (ml *migrateLogger) Printf(format string, v ...interface{}) {
	ml.logger.Info().Msg(fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...))
}

func (ml *migrateLogger) Verbose() bool {
	return true
}
