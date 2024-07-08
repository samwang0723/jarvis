package sqlc

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	esdb "github.com/samwang0723/jarvis/internal/eventsourcing/db"
)

type BalanceLoaderSaver struct {
	queries *sqlcdb.Queries
}

func (bls *BalanceLoaderSaver) Load(
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

func (bls *BalanceLoaderSaver) Save(ctx context.Context, aggregate eventsourcing.Aggregate) error {
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
