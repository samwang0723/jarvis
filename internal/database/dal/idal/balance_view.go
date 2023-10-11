package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type IBalanceViewDAL interface {
	GetBalanceViewByUserID(ctx context.Context, id uint64) (*entity.BalanceView, error)
}
