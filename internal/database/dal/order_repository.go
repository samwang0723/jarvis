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

func (i *dalImpl) CreateOrder(
	ctx context.Context,
	orderRequest *entity.Order,
	transactions []*entity.Transaction,
	orderType string,
) error {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		ctx = database.WithTx(ctx, tx)
		// check non-closed orders to perform sell or buy
		var objs []*entity.Order
		condition := ""
		if orderType == entity.OrderTypeSell {
			condition = "and buy_quantity - sell_quantity > 0"
		} else {
			condition = "and sell_quantity - buy_quantity > 0"
		}

		sql := fmt.Sprintf(`select * from orders where 
                        user_id = %d and stock_id = %d 
                        and status IN ('created', 'changed')
                        %s order by created_at asc;`, orderRequest.UserID, orderRequest.StockID, condition)

		err := i.db.Raw(sql).Scan(&objs).Error
		if err != nil {
			return err
		}
		// loop through all orders to perform sell or buy
		// if orderRequest is sell, then loop through all buy open orders until satisfied
		// if orderRequest is buy, then loop through all sell open orders until satisfied
		remainingQuantity := uint64(0)
		if orderType == entity.OrderTypeSell {
			remainingQuantity = orderRequest.SellQuantity
		} else {
			remainingQuantity = orderRequest.BuyQuantity
		}

		for _, order := range objs {
			obj, err := i.orderRepository.Load(ctx, order.ID)
			if err != nil {
				return err
			}

			// loop through all open orders until satisfied
			if orderType == entity.OrderTypeSell {
				if obj.BuyQuantity-obj.SellQuantity >= remainingQuantity {
					obj.SellPrice = (obj.SellPrice*float32(obj.SellQuantity)+orderRequest.SellPrice*float32(remainingQuantity))/float32(obj.SellQuantity) + float32(remainingQuantity)
					obj.SellQuantity += remainingQuantity
					obj.SellExchangeDate = orderRequest.SellExchangeDate
					remainingQuantity = 0
				} else if obj.BuyQuantity-obj.SellQuantity < remainingQuantity {
					obj.SellPrice = (obj.SellPrice*float32(obj.SellQuantity) + orderRequest.SellPrice*float32(obj.BuyQuantity-obj.SellQuantity)) / float32(obj.BuyQuantity)
					obj.SellQuantity = obj.BuyQuantity
					obj.SellExchangeDate = orderRequest.SellExchangeDate
					remainingQuantity -= obj.BuyQuantity - obj.SellQuantity
				}
			} else {
				if obj.SellQuantity-obj.BuyQuantity >= remainingQuantity {
					obj.BuyPrice = (obj.BuyPrice*float32(obj.BuyQuantity)+orderRequest.BuyPrice*float32(remainingQuantity))/float32(obj.BuyQuantity) + float32(remainingQuantity)
					obj.BuyQuantity += remainingQuantity
					obj.BuyExchangeDate = orderRequest.BuyExchangeDate
					remainingQuantity = 0
				} else if obj.SellQuantity-obj.BuyQuantity < remainingQuantity {
					obj.BuyPrice = (obj.BuyPrice*float32(obj.BuyQuantity) + orderRequest.BuyPrice*float32(obj.SellQuantity-obj.BuyQuantity)) / float32(obj.SellQuantity)
					obj.BuyQuantity = obj.SellQuantity
					obj.BuyExchangeDate = orderRequest.BuyExchangeDate
					remainingQuantity -= obj.SellQuantity - obj.BuyQuantity
				}
			}
			if err := i.orderRepository.Save(ctx, obj); err != nil {
				return err
			}
		}

		// No remaining pending quantity, need to create a new order
		if remainingQuantity > 0 {
			if orderType == entity.OrderTypeSell {
				orderRequest.SellQuantity = remainingQuantity
			} else {
				orderRequest.BuyQuantity = remainingQuantity
			}

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
