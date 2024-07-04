package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

type Connection struct {
	queries *sqlcdb.Queries
	pool    *pgxpool.Pool
}

func NewConnection(pool *pgxpool.Pool) *Connection {
	return &Connection{
		queries: sqlcdb.New(pool),
		pool:    pool,
	}
}

type Repo struct {
	primaryConn *Connection
	replicaConn *Connection
	logger      *zerolog.Logger
}

func NewSQLCRepository(pool *pgxpool.Pool, logger *zerolog.Logger, opts ...Option) *Repo {
	repo := &Repo{
		primaryConn: NewConnection(pool),
		logger:      logger,
	}

	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

type Option func(*Repo)

// WithReplica flexibility to choose whether to have slave pool
// Can be optimized to pass in multiple pools, if needed later
func WithReplica(pool *pgxpool.Pool) Option {
	return func(repo *Repo) {
		if pool != nil {
			repo.replicaConn = NewConnection(pool)
		}
	}
}

func (repo *Repo) primary() *sqlcdb.Queries {
	return repo.primaryConn.queries
}

//lint:ignore U1000 This is a placeholder
func (repo *Repo) replica() *sqlcdb.Queries {
	if repo.replicaConn == nil {
		return repo.primaryConn.queries
	}

	return repo.replicaConn.queries
}

//lint:ignore U1000 This is a placeholder
func (repo *Repo) primaryPgPool() *pgxpool.Pool {
	return repo.primaryConn.pool
}

//lint:ignore U1000 This is a placeholder
func (repo *Repo) replicaPgPool() *pgxpool.Pool {
	if repo.replicaConn == nil {
		return repo.primaryConn.pool
	}

	return repo.replicaConn.pool
}

func run() error {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	connString := "postgres://jarvis_app:abcd1234@localhost:5432/jarvis_main"
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	repo := NewSQLCRepository(pool, &logger, WithReplica(nil))
	stocks, err := repo.primary().ListStocks(ctx, &sqlcdb.ListStocksParams{
		Limit:           10,
		Offset:          0,
		StockIds:        []string{"1101", "2330"},
		FilterByStockID: true,
	})
	if err != nil {
		return err
	}
	// print stocks
	for _, stock := range stocks {
		log.Printf("Stock: %v\n", stock)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
