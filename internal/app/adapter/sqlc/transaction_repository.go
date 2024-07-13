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
	"github.com/samwang0723/jarvis/internal/helper"
)

type transactionLoaderSaver struct {
	queries *sqlcdb.Queries
}

func (ols *transactionLoaderSaver) Load(
	ctx context.Context,
	id uuid.UUID,
) (eventsourcing.Aggregate, error) {
	queries := ols.queries

	if trans, ok := esdb.GetTx(ctx); ok {
		queries = ols.queries.WithTx(trans)
	}

	sqlcTrans, err := queries.GetTransaction(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newRecordNotFoundError(err)
		}

		return nil, fmt.Errorf("queries.GetTransaction error: %w", err)
	}

	return fromSqlcTransaction(sqlcTrans), nil
}

func (ols *transactionLoaderSaver) Save(
	ctx context.Context,
	aggregate eventsourcing.Aggregate,
) error {
	queries := ols.queries

	if trans, ok := esdb.GetTx(ctx); ok {
		queries = ols.queries.WithTx(trans)
	}

	trans, ok := aggregate.(*domain.Transaction)
	if !ok {
		return &TypeMismatchError{
			expect: &domain.Transaction{},
			got:    aggregate,
		}
	}

	if err := queries.UpsertTransaction(ctx, &sqlcdb.UpsertTransactionParams{
		ID:           trans.ID,
		UserID:       trans.UserID,
		OrderID:      trans.OrderID,
		OrderType:    trans.OrderType,
		CreditAmount: helper.Float32ToDecimal(trans.CreditAmount),
		Status:       trans.Status,
		Version:      int32(trans.Version),
	}); err != nil {
		return fmt.Errorf("queries.UpsertTransaction error: %w", err)
	}

	return nil
}

func fromSqlcTransaction(sqlcTrans *sqlcdb.Transaction) *domain.Transaction {
	return &domain.Transaction{
		CreatedAt: sqlcTrans.CreatedAt,
		UpdatedAt: sqlcTrans.UpdatedAt,
		BaseAggregate: eventsourcing.BaseAggregate{
			ID:      sqlcTrans.ID,
			Version: int(sqlcTrans.Version),
		},
		OrderType:    sqlcTrans.OrderType,
		Status:       sqlcTrans.Status,
		UserID:       sqlcTrans.UserID,
		OrderID:      sqlcTrans.OrderID,
		CreditAmount: helper.DecimalToFloat32(sqlcTrans.CreditAmount),
		DebitAmount:  helper.DecimalToFloat32(sqlcTrans.DebitAmount),
	}
}

type transactionRepository struct {
	repo *esdb.AggregateRepository
}

func newTransactionRepository(dbPool *pgxpool.Pool) *transactionRepository {
	loaderSaver := &transactionLoaderSaver{
		queries: sqlcdb.New(dbPool),
	}

	return &transactionRepository{
		repo: esdb.NewAggregateRepository(
			&domain.Transaction{},
			dbPool,
			esdb.WithAggregateLoader(loaderSaver),
			esdb.WithAggregateSaver(loaderSaver),
		),
	}
}

func (tr *transactionRepository) Load(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Transaction, error) {
	aggregate, err := tr.repo.Load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to transactionRepository.Load: %w", err)
	}

	trans, ok := aggregate.(*domain.Transaction)
	if !ok {
		return nil, &TypeMismatchError{expect: &domain.Transaction{}, got: aggregate}
	}

	return trans, nil
}

func (tr *transactionRepository) Save(ctx context.Context, order *domain.Transaction) error {
	err := tr.repo.Save(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to transactionRepository.Save: %w", err)
	}

	return nil
}

func (repo *Repo) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	trasactions := []*domain.Transaction{transaction}
	return repo.createChainTransactions(ctx, trasactions)
}

func (repo *Repo) createChainTransactions(
	ctx context.Context,
	transactions []*domain.Transaction,
) error {
	balanceView, err := repo.balanceRepository.Load(ctx, transactions[0].UserID)
	if err != nil {
		return fmt.Errorf("failed to createChainTransactions: %w", err)
	}

	for _, transaction := range transactions {
		// immediately completed the transaction as no external vendor dependency
		if err := transaction.Complete(); err != nil {
			return err
		}
		if err := repo.transactionRepository.Save(ctx, transaction); err != nil {
			return err
		}

		if err := moveFund(balanceView, transaction); err != nil {
			return err
		}
	}

	return repo.balanceRepository.Save(ctx, balanceView)
}

func moveFund(balanceView *domain.BalanceView, transaction *domain.Transaction) error {
	switch transaction.OrderType {
	case domain.OrderTypeBuy, domain.OrderTypeFee, domain.OrderTypeTax, domain.OrderTypeWithdraw:
		if err := balanceView.MoveAvailableToPending(transaction); err != nil {
			return err
		}

		if err := balanceView.DebitPending(transaction); err != nil {
			return err
		}
	case domain.OrderTypeSell, domain.OrderTypeDeposit:
		if err := balanceView.CreditPending(transaction); err != nil {
			return err
		}

		if err := balanceView.MovePendingToAvailable(transaction); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown order type: %s", transaction.OrderType)
	}

	return nil
}
