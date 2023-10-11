package dal

import (
	"context"
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/eventsourcing"
	"github.com/samwang0723/jarvis/internal/eventsourcing/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceRepository struct {
	repo  *db.AggregateRepository
	dbRef *gorm.DB
}

type BalanceLoaderSaver struct {
	dbRef *gorm.DB
}

func (bls *BalanceLoaderSaver) Load(ctx context.Context, id uint64) (eventsourcing.Aggregate, error) {
	balanceView := &entity.BalanceView{}

	if err := bls.dbRef.Where("id = ?", id).First(balanceView).Error; err != nil {
		return nil, err
	}

	return balanceView, nil
}

func (bls *BalanceLoaderSaver) Save(ctx context.Context, aggregate eventsourcing.Aggregate) error {
	balanceView, ok := aggregate.(*entity.BalanceView)
	if !ok {
		return &TypeMismatchError{
			expect: &entity.BalanceView{},
			got:    aggregate,
		}
	}

	if err := bls.dbRef.Save(balanceView).Error; err != nil {
		return err
	}

	return nil
}

func NewBalanceRepository(dbPool *gorm.DB) *BalanceRepository {
	loaderSaver := &BalanceLoaderSaver{
		dbRef: dbPool,
	}

	return &BalanceRepository{
		repo: db.NewAggregateRepository(&entity.BalanceView{}, dbPool,
			db.WithAggregateLoader(loaderSaver), db.WithAggregateSaver(loaderSaver),
		),
		dbRef: dbPool,
	}
}

func (br *BalanceRepository) Load(ctx context.Context, id uint64) (*entity.BalanceView, error) {
	aggregate, err := br.repo.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	balanceView, ok := aggregate.(*entity.BalanceView)
	if !ok {
		return nil, &TypeMismatchError{
			expect: &entity.BalanceView{},
			got:    aggregate,
		}
	}

	return balanceView, nil
}

func (br *BalanceRepository) Save(ctx context.Context, balanceView *entity.BalanceView) error {
	err := br.repo.Save(ctx, balanceView)
	if err != nil {
		return fmt.Errorf("failed to load balance_view: %w", err)
	}

	return nil
}

func (br *BalanceRepository) LoadForUpdate(ctx context.Context, id uint64) (*entity.BalanceView, error) {
	balanceView := &entity.BalanceView{}

	if err := br.dbRef.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).First(balanceView).Error; err != nil {
		return nil, err
	}

	return balanceView, nil
}

func (i *dalImpl) GetBalanceViewByUserID(ctx context.Context, id uint64) (*entity.BalanceView, error) {
	res := &entity.BalanceView{}
	if err := i.db.First(res, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return res, nil
}
