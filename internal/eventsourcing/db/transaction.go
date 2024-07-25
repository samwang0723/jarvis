package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextTxKey struct{}

// Transaction runs txFunc inside a db transaction.
//
// It's useful when running multiple components in the same transaction. For example:
//
//	func (r *MyRepository) Save(agg MyAggregate) error {
//	   err := Transaction(ctx, r.db, func(ctx, tx) error {
//	   	querier := r.sqlcQuerier.WithTx(tx)
//	   	eventStore := r.eventStore.WithTx(tx)
//	   	...
//	   })
//	})
//
// If ctx contains a transaction already, it starts a nested transaction.
//
// If txFunc returns a non-nil error, the transaction will be rollbacked. If it's inside
// a nested transaction, only the inner transaction will be rollbacked, the outer one is
// not affeted.
func Transaction(
	ctx context.Context,
	dbPool *pgxpool.Pool,
	txFunc func(ctx context.Context, tx pgx.Tx) error,
) error {
	pgtx, ok := GetTx(ctx)

	if ok { // already inside a transaction, start nested transaction
		tx, err := pgtx.Begin(ctx) //nolint:varnamelen
		if err != nil {
			return &TransactionError{err: err, msg: "transaction error"}
		}

		defer tx.Rollback(ctx) //nolint: errcheck

		if err := txFunc(WithTx(ctx, tx), tx); err != nil {
			return fmt.Errorf("tx func error: %w", err)
		}

		return tx.Commit(ctx) //nolint:wrapcheck
	}

	// start a new transaction
	err := pgx.BeginTxFunc(ctx, dbPool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		ctx = WithTx(ctx, tx)

		return txFunc(ctx, tx)
	})
	if err != nil {
		return &TransactionError{err: err, msg: "transaction error"}
	}

	return nil
}

// GetTx retrevies a pgx.tx from context

func GetTx(ctx context.Context) (pgx.Tx, bool) { //nolint: ireturn // pgx.tx is a interface
	tx, ok := ctx.Value(contextTxKey{}).(pgx.Tx)

	return tx, ok
}

// WithTx adds a pgx.tx from context
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, contextTxKey{}, tx)
}
