package dal

import (
	"context"
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/database"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/samwang0723/jarvis/internal/eventsourcing/db"
	"gorm.io/gorm"
)

// SQL query for debugging
// SELECT aggregate_id, parent_id, event_type, version, CONVERT(payload USING utf8) as pstr FROM transaction_events;

type TransactionRepository struct {
	repo  *db.AggregateRepository
	dbRef *gorm.DB
	query *database.Query
}

type TransactionLoaderSaver struct {
	dbRef *gorm.DB
	query *database.Query
}

func (tls *TransactionLoaderSaver) Load(
	ctx context.Context,
	id uint64,
) (eventsourcing.Aggregate, error) {
	queries := tls.query

	if trans, ok := database.GetTx(ctx); ok {
		queries = tls.query.WithTx(trans)
	}

	transaction := &entity.Transaction{}
	if err := queries.Where("id = ?", id).First(transaction).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (tls *TransactionLoaderSaver) Save(
	ctx context.Context,
	aggregate eventsourcing.Aggregate,
) error {
	queries := tls.query

	if trans, ok := database.GetTx(ctx); ok {
		queries = tls.query.WithTx(trans)
	}

	transaction, ok := aggregate.(*entity.Transaction)
	if !ok {
		return &TypeMismatchError{
			expect: &entity.Transaction{},
			got:    aggregate,
		}
	}

	return queries.Omit("OriginalExchangeDate").Save(transaction).Error
}

func NewTransactionRepository(dbPool *gorm.DB) *TransactionRepository {
	loaderSaver := &TransactionLoaderSaver{
		dbRef: dbPool,
		query: database.NewQuery(dbPool),
	}

	return &TransactionRepository{
		repo: db.NewAggregateRepository(&entity.Transaction{}, dbPool,
			db.WithAggregateLoader(loaderSaver), db.WithAggregateSaver(loaderSaver),
		),
		dbRef: dbPool,
		query: database.NewQuery(dbPool),
	}
}

func (tr *TransactionRepository) Load(ctx context.Context, id uint64) (*entity.Transaction, error) {
	aggregate, err := tr.repo.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	transactionRequest, ok := aggregate.(*entity.Transaction)
	if !ok {
		return nil, &TypeMismatchError{
			got:    aggregate,
			expect: &entity.Transaction{},
		}
	}

	return transactionRequest, nil
}

func (tr *TransactionRepository) Save(
	ctx context.Context,
	transactionRequest *entity.Transaction,
) error {
	err := tr.repo.Save(ctx, transactionRequest)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

func (i *dalImpl) CreateChainTransactions(
	ctx context.Context,
	transactions []*entity.Transaction,
) error {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		ctx = database.WithTx(ctx, tx)
		balanceView, err := i.balanceRepository.LoadForUpdate(ctx, transactions[0].UserID)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			// immediately completed the transaction as no external vendor dependency
			if err := transaction.Complete(); err != nil {
				return err
			}

			if err := i.transactionRepository.Save(ctx, transaction); err != nil {
				return err
			}

			if err := moveFund(balanceView, transaction); err != nil {
				return err
			}
		}

		return i.balanceRepository.Save(ctx, balanceView)
	})

	return err
}

func moveFund(balanceView *entity.BalanceView, transaction *entity.Transaction) error {
	switch transaction.OrderType {
	case entity.OrderTypeBuy, entity.OrderTypeFee, entity.OrderTypeTax, entity.OrderTypeWithdraw:
		if err := balanceView.MoveAvailableToPending(transaction); err != nil {
			return err
		}

		if err := balanceView.DebitPending(transaction); err != nil {
			return err
		}
	case entity.OrderTypeSell, entity.OrderTypeDeposit:
		if err := balanceView.CreditPending(transaction); err != nil {
			return err
		}

		if err := balanceView.MovePendingToAvailable(transaction); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownOrderType, transaction.OrderType)
	}

	return nil
}
