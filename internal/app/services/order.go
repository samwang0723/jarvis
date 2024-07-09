package services

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	taiwanStockQuantity = 1000
	dayTradeTaxRate     = 0.5
	taxRate             = 0.003
	feeRate             = 0.001425
	brokerFeeDiscount   = 0.25
)

type processedOrder struct {
	order            *domain.Order
	exchangeQuantity uint64
}

//nolint:nolintlint,cyclop
func (s *serviceImpl) ListOrders(
	ctx context.Context,
	req *dto.ListOrderRequest,
) (objs []*domain.Order, totalCount int64, err error) {
	params := &domain.ListOrdersParams{
		UserID: s.currentUserID,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	if req.SearchParams.ExchangeMonth != nil {
		params.ExchangeMonth = *req.SearchParams.ExchangeMonth
	}
	if req.SearchParams.Status != nil {
		params.Status = *req.SearchParams.Status
	}
	if req.SearchParams.StockIDs != nil {
		params.StockIDs = *req.SearchParams.StockIDs
	}

	objs, err = s.dal.ListOrders(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	m := helper.SliceToMap(objs, func(order *domain.Order) string {
		return order.StockID
	})

	// pull lastest price from database
	stockIDs := helper.Keys(m)
	p, err := s.dal.ListLatestPrice(ctx, stockIDs)
	if err != nil {
		return nil, 0, err
	}

	stocks, _, err := s.ListStock(ctx, &dto.ListStockRequest{
		SearchParams: &dto.ListStockSearchParams{
			StockIDs: &stockIDs,
			Country:  "TW",
		},
		Offset: 0,
		Limit:  int32(len(stockIDs)),
	})
	if err != nil {
		return nil, 0, err
	}

	// calculate settled profit loss
	for _, order := range objs {
		for _, stock := range stocks {
			if stock.ID == order.StockID {
				order.StockName = stock.Name

				break
			}
		}
		order.CalculateProfitLoss()
	}

	// calculate unrealized profit loss
	for _, price := range p {
		for _, order := range objs {
			if order.StockID == price.StockID && order.Status != "closed" {
				order.CalculateUnrealizedProfitLoss(price.Price)
			}
		}
	}

	err = s.fillRealtimePrice(ctx, objs)
	if err != nil {
		return nil, 0, err
	}

	return objs, int64(len(objs)), nil
}

func (s *serviceImpl) fillRealtimePrice(ctx context.Context, objs []*domain.Order) error {
	realtimeList, err := s.fetchRealtimePrice(ctx)
	if err != nil {
		return err
	}
	for _, order := range objs {
		// override realtime data with history record.
		realtime, ok := realtimeList[order.StockID]
		if ok {
			order.CalculateUnrealizedProfitLoss(realtime.Close)
		}
	}

	return nil
}

func (s *serviceImpl) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) error {
	// check remaining open buy or sell order has quantity left to fulfill based on order type
	remainingOrders, err := s.dal.ListOpenOrders(ctx, s.currentUserID, req.StockID, req.OrderType)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list open orders")

		return err
	}

	processedOrders := []*processedOrder{}
	pendingQuantity := req.Quantity
	for _, order := range remainingOrders {
		merged, left := s.mergeOrderQuantity(order, req, pendingQuantity)
		processedOrders = append(processedOrders, &processedOrder{
			order:            order,
			exchangeQuantity: merged,
		})
		pendingQuantity = left

		if pendingQuantity == 0 {
			break
		}
	}

	// cannot fulfill based on existing open orders
	if pendingQuantity > 0 {
		order, err := domain.NewOrder(
			s.currentUserID,
			req.OrderType,
			req.StockID,
			req.ExchangeDate,
			req.TradePrice,
			pendingQuantity,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to create order")

			return err
		}
		processedOrders = append(processedOrders, &processedOrder{
			order:            order,
			exchangeQuantity: pendingQuantity,
		})
	}

	saveOrders := []*domain.Order{}
	processedTrans := []*domain.Transaction{}
	for _, po := range processedOrders {
		order := po.order
		dayTrade := order.BuyExchangeDate == order.SellExchangeDate
		partialCloseOrClose := (order.BuyQuantity > 0 && order.SellQuantity > 0)
		transactions, err := s.chainTransactions(
			order.ID,
			order.UserID,
			req.TradePrice,
			po.exchangeQuantity,
			req.OrderType,
			partialCloseOrClose,
			dayTrade,
		)
		if err != nil {
			return errUnableToChainTransactions
		}
		processedTrans = append(processedTrans, transactions...)
		saveOrders = append(saveOrders, order)
	}

	return s.dal.CreateOrder(ctx, saveOrders, processedTrans)
}

func (s *serviceImpl) mergeOrderQuantity(
	order *domain.Order,
	req *dto.CreateOrderRequest,
	pendingQuantity uint64,
) (mergedQuantity, leftQuantity uint64) {
	// if buy = 4, sell = 0, quantity = 2, orderType = sell, then sell = 2
	// if buy = 4, sell = 3, quantity = 2, orderType = sell, then sell = 4, left = 1
	leftQuantity = pendingQuantity
	eventQuantity := uint64(0)
	price := float32(0.0)

	switch req.OrderType {
	case domain.OrderTypeBuy:
		eventQuantity = order.BuyQuantity
		gap := order.SellQuantity - order.BuyQuantity
		if gap >= leftQuantity {
			price = (order.BuyPrice*float32(order.BuyQuantity) + req.TradePrice*float32(req.Quantity)) / float32(
				order.BuyQuantity+req.Quantity,
			)
			mergedQuantity = leftQuantity
			leftQuantity = 0
		} else {
			price = (order.BuyPrice*float32(order.BuyQuantity) + req.TradePrice*float32(gap)) / float32(order.BuyQuantity+gap)
			mergedQuantity = gap
			leftQuantity -= gap
		}
		eventQuantity += mergedQuantity
	case domain.OrderTypeSell:
		eventQuantity = order.SellQuantity
		gap := order.BuyQuantity - order.SellQuantity
		if gap >= leftQuantity {
			price = (order.SellPrice*float32(order.SellQuantity) + req.TradePrice*float32(req.Quantity)) / float32(
				order.SellQuantity+req.Quantity,
			)
			mergedQuantity = leftQuantity
			leftQuantity = 0
		} else {
			price = (order.SellPrice*float32(order.SellQuantity) + req.TradePrice*float32(gap)) / float32(order.SellQuantity+gap)
			mergedQuantity = gap
			leftQuantity -= gap
		}
		eventQuantity += mergedQuantity
	}

	err := order.Change(
		req.OrderType,
		req.StockID,
		req.ExchangeDate,
		price,
		eventQuantity,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to change order")
	}

	return mergedQuantity, leftQuantity
}

func (s *serviceImpl) chainTransactions(
	orderID uuid.UUID,
	userID uuid.UUID,
	price float32,
	quantity uint64,
	orderType string,
	partialCloseOrClose bool,
	dayTrade bool,
) (chainedTransactions []*domain.Transaction, err error) {
	debitAmount, creditAmount := float32(0.0), float32(0.0)
	switch orderType {
	case domain.OrderTypeBuy:
		debitAmount = price * float32(quantity) * taiwanStockQuantity
	case domain.OrderTypeSell:
		creditAmount = price * float32(quantity) * taiwanStockQuantity
	}

	transaction, err := domain.NewTransaction(
		userID,
		orderType,
		creditAmount,
		debitAmount,
		orderID,
	)
	if err != nil {
		return chainedTransactions, err
	}

	chainedTransactions = append(chainedTransactions, transaction)

	tax, err := s.genTaxTransaction(
		orderID,
		userID,
		price,
		quantity,
		orderType,
		partialCloseOrClose,
		dayTrade,
	)
	if err != nil {
		return chainedTransactions, err
	} else if tax != nil {
		chainedTransactions = append(chainedTransactions, tax)
	}

	fee, err := s.genFeeTransaction(orderID, userID, price, quantity, orderType)
	if err != nil {
		return chainedTransactions, err
	} else if fee != nil {
		chainedTransactions = append(chainedTransactions, fee)
	}

	return chainedTransactions, nil
}

func (s *serviceImpl) genTaxTransaction(
	orderID uuid.UUID,
	userID uuid.UUID,
	price float32,
	quantity uint64,
	orderType string,
	partialCloseOrClose bool,
	dayTrade bool,
) (*domain.Transaction, error) {
	// only charge tax on partial order close or complete order close
	if partialCloseOrClose {
		debitAmount := float32(0.0)
		if orderType == domain.OrderTypeBuy {
			debitAmount = price * float32(quantity) * taiwanStockQuantity * taxRate
		} else if orderType == domain.OrderTypeSell {
			debitAmount = price * float32(quantity) * taiwanStockQuantity * taxRate
		}
		if dayTrade {
			debitAmount *= dayTradeTaxRate
		}

		output, err := domain.NewTransaction(
			userID,
			domain.OrderTypeTax,
			0,
			debitAmount,
			orderID,
		)

		return output, err
	}

	//nolint:nilnil // this is a special case
	return nil, nil
}

func (s *serviceImpl) genFeeTransaction(
	orderID uuid.UUID,
	userID uuid.UUID,
	price float32,
	quantity uint64,
	orderType string,
) (*domain.Transaction, error) {
	debitAmount := float32(0.0)
	if orderType == domain.OrderTypeBuy {
		debitAmount = price * float32(quantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	} else if orderType == domain.OrderTypeSell {
		debitAmount = price * float32(quantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	}

	output, err := domain.NewTransaction(
		userID,
		domain.OrderTypeFee,
		0,
		debitAmount,
		orderID,
	)

	return output, err
}
