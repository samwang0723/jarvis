package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type ITransactionDAL interface {
	CreateChainTransactions(ctx context.Context, transactions []*entity.Transaction) error
}
