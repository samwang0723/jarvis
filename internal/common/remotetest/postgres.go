package remotetest

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	"github.com/samwang0723/jarvis/database"
	"github.com/samwang0723/jarvis/internal/common"
	"github.com/samwang0723/jarvis/internal/db/pginit"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // because this package is used for testing
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const (
	maxIdleConns = 10
	maxConns     = 10
)

type PostgresContainer struct {
	cfg *pginit.Config
}

func CreatePostgres() (*PostgresContainer, error) {
	host, port := "localhost", "5432"
	username, password := "postgres", "postgres"

	if envHost, got := os.LookupEnv("TEST_DATABASE_HOST"); got {
		host = envHost
	}

	if envPort, got := os.LookupEnv("TEST_DATABASE_PORT"); got {
		port = envPort
	}

	if envUsername, got := os.LookupEnv("TEST_DATABASE_USERNAME"); got {
		username = envUsername
	}

	if envPass, got := os.LookupEnv("TEST_DATABASE_PASSWORD"); got {
		password = envPass
	}

	return &PostgresContainer{
		cfg: &pginit.Config{
			User:         username,
			Password:     password,
			Host:         host,
			Port:         port,
			Database:     newDBName(),
			MaxConns:     maxConns,
			MaxIdleConns: maxIdleConns,
			MaxLifeTime:  time.Minute,
		},
	}, nil
}

func newDBName() string {
	validDBID := strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "_")

	return fmt.Sprintf("knock_out_test_%s", validDBID)
}

func (pg *PostgresContainer) CreateConnPool() (*pgxpool.Pool, error) {
	logger := zerolog.New(io.Discard)

	if err := pg.createDB(); err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	pgi, err := pginit.New(pg.cfg,
		pginit.WithLogLevel(zerolog.WarnLevel),
		pginit.WithLogger(&logger, "request-id"),
		pginit.WithUUIDType(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not init pginit: %w", err)
	}

	var connPool *pgxpool.Pool

	if err = common.ExponentialBackoffRetry(func() error {
		var poolErr error
		connPool, poolErr = pgi.ConnPool(context.Background())

		if poolErr != nil {
			return fmt.Errorf("pgi connection error: %w", poolErr)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return connPool, nil
}

func (pg *PostgresContainer) getConnectionURL(dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		pg.cfg.User, pg.cfg.Password, net.JoinHostPort(pg.cfg.Host, pg.cfg.Port),
		dbName)
}

func (pg *PostgresContainer) createDB() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, pg.getConnectionURL(""))
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	defer conn.Close(ctx)

	// somehow CREATE DATABASE $1; doesn't work
	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s;", pg.cfg.Database))
	if err != nil {
		return fmt.Errorf("failed to create db: %w", err)
	}

	return nil
}

func (pg *PostgresContainer) RunMigrations() error {
	driver, err := iofs.New(database.MigrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("iofs error: %w", err)
	}

	defer driver.Close()

	dbURL := pg.getConnectionURL(pg.cfg.Database)

	migrator, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		return fmt.Errorf("migrate new error: %w", err)
	}

	defer migrator.Close()

	err = migrator.Up()
	if err != nil {
		return fmt.Errorf("migrate up error: %w", err)
	}

	return nil
}

func (pg *PostgresContainer) Purge() error {
	ctx := context.Background()
	dbURL := pg.getConnectionURL("")

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	defer conn.Close(ctx)

	// force drop cuz connections in pool might not be closed completely yet
	dropSQL := fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", pg.cfg.Database)

	_, err = conn.Exec(ctx, dropSQL)
	if err != nil {
		return fmt.Errorf("failed to cleanup database: %w", err)
	}

	return nil
}

func SetupPostgresClient(t *testing.T, runMigration bool) *pgxpool.Pool {
	t.Helper()

	server, err := CreatePostgres()
	if err != nil {
		t.Fatalf("failed to create postgres server: %v", err)
	}

	t.Cleanup(func() {
		if errP := server.Purge(); errP != nil {
			t.Fatalf("failed to purge postgres server: %v", errP)
		}
	})

	client, err := server.CreateConnPool()
	if err != nil {
		t.Fatalf("failed to create connection pool: %v", err)
	}

	t.Cleanup(client.Close)

	if runMigration {
		if err := server.RunMigrations(); err != nil {
			t.Fatalf("failed to run migrations: %v", err)
		}
	}

	return client
}
