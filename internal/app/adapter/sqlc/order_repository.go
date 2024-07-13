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

type orderLoaderSaver struct {
	queries *sqlcdb.Queries
}

func (ols *orderLoaderSaver) Load(
	ctx context.Context,
	id uuid.UUID,
) (eventsourcing.Aggregate, error) {
	queries := ols.queries

	if trans, ok := esdb.GetTx(ctx); ok {
		queries = ols.queries.WithTx(trans)
	}

	sqlcOrder, err := queries.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, newRecordNotFoundError(err)
		}

		return nil, fmt.Errorf("queries.GetOrderView error: %w", err)
	}

	return fromSqlcOrderView(sqlcOrder), nil
}

func (ols *orderLoaderSaver) Save(ctx context.Context, aggregate eventsourcing.Aggregate) error {
	queries := ols.queries

	if trans, ok := esdb.GetTx(ctx); ok {
		queries = ols.queries.WithTx(trans)
	}

	order, ok := aggregate.(*domain.Order)
	if !ok {
		return &TypeMismatchError{
			expect: &domain.Order{},
			got:    aggregate,
		}
	}

	if err := queries.UpsertOrder(ctx, &sqlcdb.UpsertOrderParams{
		ID:               order.ID,
		UserID:           order.UserID,
		StockID:          order.StockID,
		BuyPrice:         helper.Float32ToDecimal(order.BuyPrice),
		BuyQuantity:      int64(order.BuyQuantity),
		BuyExchangeDate:  order.BuyExchangeDate,
		SellPrice:        helper.Float32ToDecimal(order.SellPrice),
		SellQuantity:     int64(order.SellQuantity),
		SellExchangeDate: order.SellExchangeDate,
		ProfitablePrice:  helper.Float32ToDecimal(order.ProfitablePrice),
		Status:           order.Status,
		Version:          int32(order.Version),
	}); err != nil {
		return fmt.Errorf("queries.UpsertOrderView error: %w", err)
	}

	return nil
}

func fromSqlcOrderView(sqlcOrder *sqlcdb.Order) *domain.Order {
	return &domain.Order{
		CreatedAt: sqlcOrder.CreatedAt,
		UpdatedAt: sqlcOrder.UpdatedAt,
		BaseAggregate: eventsourcing.BaseAggregate{
			ID:      sqlcOrder.ID,
			Version: int(sqlcOrder.Version),
		},
		StockID:          sqlcOrder.StockID,
		BuyExchangeDate:  sqlcOrder.BuyExchangeDate,
		Status:           sqlcOrder.Status,
		SellExchangeDate: sqlcOrder.SellExchangeDate,
		SellQuantity:     uint64(sqlcOrder.SellQuantity),
		BuyQuantity:      uint64(sqlcOrder.BuyQuantity),
		UserID:           sqlcOrder.UserID,
		ProfitablePrice:  helper.DecimalToFloat32(sqlcOrder.ProfitablePrice),
		SellPrice:        helper.DecimalToFloat32(sqlcOrder.SellPrice),
		BuyPrice:         helper.DecimalToFloat32(sqlcOrder.BuyPrice),
	}
}

type orderRepository struct {
	repo *esdb.AggregateRepository
}

func newOrderRepository(dbPool *pgxpool.Pool) *orderRepository {
	loaderSaver := &orderLoaderSaver{
		queries: sqlcdb.New(dbPool),
	}

	return &orderRepository{
		repo: esdb.NewAggregateRepository(
			&domain.Order{},
			dbPool,
			esdb.WithAggregateLoader(loaderSaver),
			esdb.WithAggregateSaver(loaderSaver),
		),
	}
}

func (or *orderRepository) Load(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	aggregate, err := or.repo.Load(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to orderRepository.Load: %w", err)
	}

	order, ok := aggregate.(*domain.Order)
	if !ok {
		return nil, &TypeMismatchError{expect: &domain.Order{}, got: aggregate}
	}

	return order, nil
}

func (or *orderRepository) Save(ctx context.Context, order *domain.Order) error {
	err := or.repo.Save(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to orderRepository.Save: %w", err)
	}

	return nil
}

func (repo *Repo) ListOrders(
	ctx context.Context,
	arg *domain.ListOrdersParams,
) ([]*domain.Order, error) {
	params := &sqlcdb.ListOrdersParams{
		UserID:        arg.UserID,
		Limit:         arg.Limit,
		Offset:        arg.Offset,
		Status:        arg.Status,
		ExchangeMonth: arg.ExchangeMonth,
	}
	if len(arg.StockIDs) > 0 {
		params.StockIds = arg.StockIDs
		params.FilterByStockID = true
	}

	rows, err := repo.primary().ListOrders(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to ListOrders: %w", err)
	}

	objs := make([]*domain.Order, len(rows))
	for idx, order := range rows {
		obj, err := repo.orderRepository.Load(ctx, order.ID)
		if err != nil {
			return objs, err
		}
		objs[idx] = obj
	}

	return objs, nil
}

func (repo *Repo) ListOpenOrders(
	ctx context.Context,
	userID uuid.UUID,
	stockID string,
	orderType string,
) ([]*domain.Order, error) {
	rows, err := repo.primary().ListOpenOrders(ctx, &sqlcdb.ListOpenOrdersParams{
		UserID:    userID,
		StockID:   stockID,
		OrderType: orderType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ListOpenOrders: %w", err)
	}

	objs := make([]*domain.Order, len(rows))
	for idx, orderID := range rows {
		obj, err := repo.orderRepository.Load(ctx, orderID)
		if err != nil {
			return objs, err
		}
		objs[idx] = obj
	}

	return objs, nil
}

func (repo *Repo) CreateOrder(
	ctx context.Context,
	orders []*domain.Order,
	transactions []*domain.Transaction,
) error {
	err := repo.RunInTransaction(ctx, func(ctx context.Context) error {
		for _, order := range orders {
			if err := repo.orderRepository.Save(ctx, order); err != nil {
				return fmt.Errorf("failed to orderRepository.Save: %w", err)
			}
		}

		return repo.createChainTransactions(ctx, transactions)
	})

	return err
}
