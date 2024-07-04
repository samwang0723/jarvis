package sqlc

import (
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

func NewSqlcRepository(pool *pgxpool.Pool, logger *zerolog.Logger, opts ...Option) *Repo {
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
