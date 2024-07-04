// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: stocks.sql

package sqlcdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const BatchUpsertStocks = `-- name: BatchUpsertStocks :exec
INSERT INTO stocks (id, name, country, category, market)
VALUES (
    unnest($1::text[]), 
    unnest($2::text[]), 
    unnest($3::text[]), 
    unnest($4::text[]),
    unnest($5::text[])
)
ON CONFLICT (id) DO UPDATE 
SET 
    name = EXCLUDED.name,
    country = EXCLUDED.country,
    category = EXCLUDED.category,
    market = EXCLUDED.market
`

type BatchUpsertStocksParams struct {
	Column1 []string
	Column2 []string
	Column3 []string
	Column4 []string
	Column5 []string
}

func (q *Queries) BatchUpsertStocks(ctx context.Context, arg *BatchUpsertStocksParams) error {
	_, err := q.db.Exec(ctx, BatchUpsertStocks,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
		arg.Column5,
	)
	return err
}

const CountStocks = `-- name: CountStocks :one
SELECT COUNT(*) FROM stocks
WHERE
    ($1::VARCHAR = '' OR country = $1)
    AND (id = ANY($2::text[]) OR NOT $3::bool)
    AND ($4::VARCHAR = '' OR name ILIKE '%' || $4 || '%')
    AND ($5::VARCHAR = '' OR category = $5)
`

type CountStocksParams struct {
	Country         string
	StockIds        []string
	FilterByStockID bool
	Name            string
	Category        string
}

func (q *Queries) CountStocks(ctx context.Context, arg *CountStocksParams) (int64, error) {
	row := q.db.QueryRow(ctx, CountStocks,
		arg.Country,
		arg.StockIds,
		arg.FilterByStockID,
		arg.Name,
		arg.Category,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const CreateStock = `-- name: CreateStock :exec
INSERT INTO stocks (id, name, country, category, market)
VALUES ($1, $2, $3, $4, $5)
`

type CreateStockParams struct {
	ID       string
	Name     string
	Country  string
	Category *string
	Market   *string
}

func (q *Queries) CreateStock(ctx context.Context, arg *CreateStockParams) error {
	_, err := q.db.Exec(ctx, CreateStock,
		arg.ID,
		arg.Name,
		arg.Country,
		arg.Category,
		arg.Market,
	)
	return err
}

const DeleteStockbyID = `-- name: DeleteStockbyID :exec
UPDATE stocks SET deleted_at = NOW() WHERE id = $1
`

func (q *Queries) DeleteStockbyID(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, DeleteStockbyID, id)
	return err
}

const DeleteStockbyStockID = `-- name: DeleteStockbyStockID :one
SELECT id, name, country, category, market, created_at, updated_at, deleted_at
FROM stocks
WHERE id = $1 AND deleted_at IS NULL LIMIT 1
`

type DeleteStockbyStockIDRow struct {
	ID        string
	Name      string
	Country   string
	Category  *string
	Market    *string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

func (q *Queries) DeleteStockbyStockID(ctx context.Context, id string) (*DeleteStockbyStockIDRow, error) {
	row := q.db.QueryRow(ctx, DeleteStockbyStockID, id)
	var i DeleteStockbyStockIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Country,
		&i.Category,
		&i.Market,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const ListCategories = `-- name: ListCategories :many
SELECT DISTINCT category FROM stocks
`

func (q *Queries) ListCategories(ctx context.Context) ([]*string, error) {
	rows, err := q.db.Query(ctx, ListCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*string
	for rows.Next() {
		var category *string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		items = append(items, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ListStocks = `-- name: ListStocks :many
SELECT id, name, country, site, category, market, created_at, updated_at, deleted_at FROM stocks
WHERE
    ($3::VARCHAR = '' OR country = $3)
    AND (id = ANY($4::text[]) OR NOT $5::bool)
    AND ($6::VARCHAR = '' OR name ILIKE '%' || $6 || '%')
    AND ($7::VARCHAR = '' OR category = $7)
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListStocksParams struct {
	Limit           int32
	Offset          int32
	Country         string
	StockIds        []string
	FilterByStockID bool
	Name            string
	Category        string
}

func (q *Queries) ListStocks(ctx context.Context, arg *ListStocksParams) ([]*Stock, error) {
	rows, err := q.db.Query(ctx, ListStocks,
		arg.Limit,
		arg.Offset,
		arg.Country,
		arg.StockIds,
		arg.FilterByStockID,
		arg.Name,
		arg.Category,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Stock
	for rows.Next() {
		var i Stock
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Country,
			&i.Site,
			&i.Category,
			&i.Market,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateStock = `-- name: UpdateStock :exec
UPDATE stocks
SET 
    name = $2,
    country = $3,
    category = $4,
    market = $5
WHERE id = $1
`

type UpdateStockParams struct {
	ID       string
	Name     string
	Country  string
	Category *string
	Market   *string
}

func (q *Queries) UpdateStock(ctx context.Context, arg *UpdateStockParams) error {
	_, err := q.db.Exec(ctx, UpdateStock,
		arg.ID,
		arg.Name,
		arg.Country,
		arg.Category,
		arg.Market,
	)
	return err
}
