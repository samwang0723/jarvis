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
// SELECT aggregate_id, parent_id, event_type, version, CONVERT(payload USING utf8) as pstr FROM order_events;

type OrderRepository struct {
	repo  *db.AggregateRepository
	dbRef *gorm.DB
	query *database.Query
}

type OrderLoaderSaver struct {
	dbRef *gorm.DB
	query *database.Query
}

func (ols *OrderLoaderSaver) Load(ctx context.Context, id uint64) (eventsourcing.Aggregate, error) {
	queries := ols.query

	if trans, ok := database.GetTx(ctx); ok {
		queries = ols.query.WithTx(trans)
	}

	order := &entity.Order{}
	if err := queries.Where("id = ?", id).First(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (ols *OrderLoaderSaver) Save(ctx context.Context, aggregate eventsourcing.Aggregate) error {
	queries := ols.query

	if trans, ok := database.GetTx(ctx); ok {
		queries = ols.query.WithTx(trans)
	}

	order, ok := aggregate.(*entity.Order)
	if !ok {
		return &TypeMismatchError{
			expect: &entity.Order{},
			got:    aggregate,
		}
	}

	if err := queries.Save(order).Error; err != nil {
		return err
	}

	return nil
}

func NewOrderRepository(dbPool *gorm.DB) *OrderRepository {
	loaderSaver := &OrderLoaderSaver{
		dbRef: dbPool,
		query: database.NewQuery(dbPool),
	}

	return &OrderRepository{
		repo: db.NewAggregateRepository(&entity.Order{}, dbPool,
			db.WithAggregateLoader(loaderSaver), db.WithAggregateSaver(loaderSaver),
		),
		dbRef: dbPool,
		query: database.NewQuery(dbPool),
	}
}

func (tr *OrderRepository) Load(ctx context.Context, id uint64) (*entity.Order, error) {
	aggregate, err := tr.repo.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	orderRequest, ok := aggregate.(*entity.Order)
	if !ok {
		return nil, &TypeMismatchError{
			got:    aggregate,
			expect: &entity.Order{},
		}
	}

	return orderRequest, nil
}

func (tr *OrderRepository) Save(ctx context.Context, orderRequest *entity.Order) error {
	err := tr.repo.Save(ctx, orderRequest)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

func (i *dalImpl) ListOpenOrders(
	ctx context.Context,
	userID uint64,
	stockID, orderType string,
) ([]*entity.Order, error) {
	var ids []uint64

	condition := ""
	if orderType == entity.OrderTypeSell {
		condition = "and buy_quantity - sell_quantity > 0"
	} else {
		condition = "and sell_quantity - buy_quantity > 0"
	}

	sql := fmt.Sprintf(`select id from orders where 
                        user_id = %d and stock_id = %s 
                        and status IN ('created', 'changed')
                        %s order by created_at asc;`, userID, stockID, condition)

	err := i.db.Raw(sql).Scan(&ids).Error
	if err != nil {
		return []*entity.Order{}, err
	}

	objs := make([]*entity.Order, len(ids))
	for idx, orderID := range ids {
		obj, err := i.orderRepository.Load(ctx, orderID)
		if err != nil {
			return objs, err
		}
		objs[idx] = obj
	}

	return objs, nil
}

func (i *dalImpl) CreateOrder(
	ctx context.Context,
	orders []*entity.Order,
	transactions []*entity.Transaction,
) error {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		ctx = database.WithTx(ctx, tx)

		for _, orderRequest := range orders {
			if err := i.orderRepository.Save(ctx, orderRequest); err != nil {
				return err
			}
		}

		if err := i.CreateChainTransactions(ctx, transactions); err != nil {
			return err
		}

		return nil
	})

	return err
}
