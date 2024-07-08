package sqlc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/domain"
	esdb "github.com/samwang0723/jarvis/internal/eventsourcing/db"

	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

type Trans struct {
	pgx.Tx
}

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

type balanceRepository struct {
	repo *esdb.AggregateRepository
}

func NewBalanceRepository(dbPool *pgxpool.Pool) *balanceRepository {
	loaderSaver := &BalanceLoaderSaver{
		queries: sqlcdb.New(dbPool),
	}

	return &balanceRepository{
		repo: esdb.NewAggregateRepository(
			&domain.BalanceView{},
			dbPool,
			esdb.WithAggregateLoader(loaderSaver),
			esdb.WithAggregateSaver(loaderSaver),
		),
	}
}

func (repo *Repo) Transaction(ctx context.Context) (*Trans, error) {
	trans, err := repo.primaryConn.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &Trans{
		Tx: trans,
	}, nil
}

func (repo *Repo) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := repo.primaryConn.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	repoTrans := &Trans{
		Tx: tx,
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			zerolog.Ctx(ctx).Error().Err(rollbackErr).Msg("failed to rollback transaction")

			err = rollbackErr
		}
	}()

	ctx = esdb.WithTx(ctx, repoTrans)
	if err := fn(ctx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
