// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: balance_view.sql

package sqlcdb

import (
	"context"

	uuid "github.com/gofrs/uuid/v5"
)

const GetBalanceView = `-- name: GetBalanceView :one
SELECT id, balance, available, pending, version, created_at, updated_at
FROM balance_views
WHERE id = $1
`

func (q *Queries) GetBalanceView(ctx context.Context, id uuid.UUID) (*BalanceView, error) {
	row := q.db.QueryRow(ctx, GetBalanceView, id)
	var i BalanceView
	err := row.Scan(
		&i.ID,
		&i.Balance,
		&i.Available,
		&i.Pending,
		&i.Version,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const UpsertBalanceView = `-- name: UpsertBalanceView :exec
INSERT INTO balance_views (id, balance, available, pending, version)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE
SET balance = EXCLUDED.balance,
    available = EXCLUDED.available,
    pending = EXCLUDED.pending,
    version = EXCLUDED.version
`

type UpsertBalanceViewParams struct {
	ID        uuid.UUID
	Balance   float64
	Available float64
	Pending   float64
	Version   int32
}

func (q *Queries) UpsertBalanceView(ctx context.Context, arg *UpsertBalanceViewParams) error {
	_, err := q.db.Exec(ctx, UpsertBalanceView,
		arg.ID,
		arg.Balance,
		arg.Available,
		arg.Pending,
		arg.Version,
	)
	return err
}