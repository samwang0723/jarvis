package sqlc

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	esdb "github.com/samwang0723/jarvis/internal/eventsourcing/db"
)

type balanceLoaderSaver struct {
	queries *sqlcdb.Queries
}

func (bls *balanceLoaderSaver) Load(
	ctx context.Context,
	id uuid.UUID,
) (eventsourcing.Aggregate, error) {
	queries := bls.queries

	if trans, ok := esdb.GetTx(ctx); ok {
		queries = bls.queries.WithTx(trans)
	}

	sqlcBalance, err := queries.GetBalanceView(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newRecordNotFoundError(err)
		}

		return nil, fmt.Errorf("queries.GetBalanceView error: %w", err)
	}

	return fromSqlcBalanceView(sqlcBalance), nil
}

func (bls *balanceLoaderSaver) Save(ctx context.Context, aggregate eventsourcing.Aggregate) error {
	queries := bls.queries

	if trans, ok := esdb.GetTx(ctx); ok {
		queries = bls.queries.WithTx(trans)
	}

	balanceView, ok := aggregate.(*domain.BalanceView)
	if !ok {
		return &TypeMismatchError{
			expect: &domain.BalanceView{},
			got:    aggregate,
		}
	}

	if err := queries.UpsertBalanceView(ctx, &sqlcdb.UpsertBalanceViewParams{
		ID:        balanceView.ID,
		Balance:   float64(balanceView.Balance),
		Available: float64(balanceView.Available),
		Pending:   float64(balanceView.Pending),
		Version:   int32(balanceView.Version),
	}); err != nil {
		return fmt.Errorf("queries.UpsertBalanceView error: %w", err)
	}

	return nil
}

func fromSqlcBalanceView(sqlcBalance *sqlcdb.BalanceView) *domain.BalanceView {
	return &domain.BalanceView{
		CreatedAt: sqlcBalance.CreatedAt.Time,
		UpdatedAt: sqlcBalance.UpdatedAt.Time,
		BaseAggregate: eventsourcing.BaseAggregate{
			ID:      sqlcBalance.ID,
			Version: int(sqlcBalance.Version),
		},
		Balance:   float32(sqlcBalance.Balance),
		Pending:   float32(sqlcBalance.Pending),
		Available: float32(sqlcBalance.Available),
	}
}

type balanceRepository struct {
	repo *esdb.AggregateRepository
}

func newBalanceRepository(dbPool *pgxpool.Pool) *balanceRepository {
	loaderSaver := &balanceLoaderSaver{
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

func (br *balanceRepository) Load(ctx context.Context, id uuid.UUID) (*domain.BalanceView, error) {
	aggregate, err := br.repo.Load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to balanceRepository.Load: %w", err)
	}

	balanceView, ok := aggregate.(*domain.BalanceView)
	if !ok {
		return nil, &TypeMismatchError{expect: &domain.BalanceView{}, got: aggregate}
	}

	return balanceView, nil
}

func (br *balanceRepository) Save(ctx context.Context, balanceView *domain.BalanceView) error {
	err := br.repo.Save(ctx, balanceView)
	if err != nil {
		return fmt.Errorf("failed to balanceRepository.Save: %w", err)
	}

	return nil
}

func (repo *Repo) GetBalanceView(ctx context.Context, id uuid.UUID) (*domain.BalanceView, error) {
	return repo.balanceRepository.Load(ctx, id)
}

func (repo *Repo) createBalance(
	ctx context.Context,
	userID uuid.UUID,
	initBalance float32,
) error {
	balanceView, err := domain.NewBalanceView(userID, initBalance)
	if err != nil {
		return fmt.Errorf("failed to apply event to balanceView: %w", err)
	}

	return repo.balanceRepository.Save(ctx, balanceView)
}
