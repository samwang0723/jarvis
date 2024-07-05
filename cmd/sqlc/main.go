package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/samwang0723/jarvis/internal/app/adapter"
	"github.com/samwang0723/jarvis/internal/app/adapter/sqlc"
	"github.com/samwang0723/jarvis/internal/app/domain"
)

func run() error {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	connString := "postgres://jarvis_app:abcd1234@localhost:5432/jarvis_main"
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to parse connection string")
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to create connection pool")
		return err
	}

	defer pool.Close()

	repo := sqlc.NewSqlcRepository(pool, &logger)
	adapter := adapter.NewAdapterImp(repo)

	categories, err := adapter.ListCategories(ctx)
	if err != nil {
		return err
	}
	// print all categories
	for _, category := range categories {
		logger.Info().Str("category", *category).Msg("Category")
	}

	stocks, err := adapter.ListStocks(ctx, &domain.ListStocksParams{
		Limit:           10,
		Offset:          0,
		Country:         "TW",
		StockIDs:        []string{"1101", "2330"},
		FilterByStockID: true,
	})
	if err != nil {
		return err
	}
	// print all stocks
	for _, stock := range stocks {
		logger.Info().Str("name", stock.Name).Str("id", stock.ID).Msg("Stock")
	}

	threePrimary, err := adapter.ListThreePrimary(ctx, &domain.ListThreePrimaryParams{
		Limit:     10,
		Offset:    0,
		StockID:   "1101",
		StartDate: "2024-01-01",
		EndDate:   "2024-02-12",
	})
	if err != nil {
		return err
	}
	// print all three primary
	for _, three := range threePrimary {
		var dealerTradeSharesStr string
		if three.DealerTradeShares != nil {
			dealerTradeSharesStr = strconv.FormatInt(*three.DealerTradeShares, 10)
		} else {
			dealerTradeSharesStr = ""
		}
		logger.Info().
			Str("stock_id", three.StockID).
			Str("exchange_date", three.ExchangeDate).
			Str("foreign_trade_shares", dealerTradeSharesStr).
			Msg("ThreePrimary")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
