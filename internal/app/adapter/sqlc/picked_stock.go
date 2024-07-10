package sqlc

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

func (repo *Repo) CreatePickedStocks(ctx context.Context, objs []*domain.PickedStock) error {
	userIDs := make([]uuid.UUID, 0, len(objs))
	stockIDs := make([]string, 0, len(objs))

	for _, obj := range objs {
		userIDs = append(userIDs, obj.UserID)
		stockIDs = append(stockIDs, obj.StockID)
	}

	return repo.primary().CreatePickedStocks(ctx, &sqlcdb.CreatePickedStocksParams{
		UserIds:  userIDs,
		StockIds: stockIDs,
	})
}

func (repo *Repo) DeletePickedStock(
	ctx context.Context,
	userID uuid.UUID,
	stockID string,
) error {
	return repo.primary().DeletePickedStock(ctx, &sqlcdb.DeletePickedStockParams{
		UserID:  userID,
		StockID: stockID,
	})
}

func (repo *Repo) ListPickedStocks(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.PickedStock, error) {
	pickedStocks, err := repo.primary().ListPickedStocks(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.PickedStock, 0, len(pickedStocks))
	for _, pickedStock := range pickedStocks {
		time := domain.Time{
			CreatedAt: &pickedStock.CreatedAt.Time,
			UpdatedAt: &pickedStock.UpdatedAt.Time,
		}
		if pickedStock.DeletedAt.Valid {
			time.DeletedAt = &pickedStock.DeletedAt.Time
		}
		result = append(result, domain.PickedStock{
			ID:      domain.ID{ID: pickedStock.ID},
			UserID:  pickedStock.UserID,
			StockID: pickedStock.StockID,
			Time:    time,
		})
	}
	return result, nil
}
