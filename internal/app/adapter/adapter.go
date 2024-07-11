package adapter

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/adapter/sqlc"
	"github.com/samwang0723/jarvis/internal/app/domain"
)

type Adapter interface {
	BatchUpsertStocks(ctx context.Context, objs []*domain.Stock) error
	CreateStock(ctx context.Context, obj *domain.Stock) error
	DeleteStockByID(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]*string, error)
	ListStocks(ctx context.Context, arg *domain.ListStocksParams) ([]*domain.Stock, error)
	BatchUpsertThreePrimary(ctx context.Context, objs []*domain.ThreePrimary) error
	CreateThreePrimary(ctx context.Context, arg *domain.ThreePrimary) error
	ListThreePrimary(
		ctx context.Context,
		arg *domain.ListThreePrimaryParams,
	) ([]*domain.ThreePrimary, error)
	BatchUpsertDailyClose(ctx context.Context, objs []*domain.DailyClose) error
	CreateDailyClose(ctx context.Context, obj *domain.DailyClose) error
	HasDailyClose(ctx context.Context, date string) (bool, error)
	ListDailyClose(
		ctx context.Context,
		arg *domain.ListDailyCloseParams,
	) ([]*domain.DailyClose, error)
	ListLatestPrice(ctx context.Context, stockIDs []string) ([]*domain.StockPrice, error)
	BatchUpsertStakeConcentration(ctx context.Context, objs []*domain.StakeConcentration) error
	GetStakeConcentrationByStockID(
		ctx context.Context,
		stockID, date string,
	) (*domain.StakeConcentration, error)
	GetStakeConcentrationsWithVolumes(
		ctx context.Context,
		stockID, date string,
	) ([]*domain.CalculationBase, error)
	HasStakeConcentration(ctx context.Context, exchangeDate string) (bool, error)
	GetStakeConcentrationLatestDataPoint(ctx context.Context) (string, error)
	CreatePickedStocks(ctx context.Context, objs []*domain.PickedStock) error
	DeletePickedStock(ctx context.Context, userID uuid.UUID, stockID string) error
	ListPickedStocks(ctx context.Context, userID uuid.UUID) ([]domain.PickedStock, error)
	CreateUser(ctx context.Context, obj *domain.User) error
	UpdateUser(ctx context.Context, obj *domain.User) error
	UpdateSessionID(ctx context.Context, params *domain.UpdateSessionIDParams) error
	DeleteSessionID(ctx context.Context, userID uuid.UUID) error
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	ListUsers(ctx context.Context, limit int32, offset int32) ([]*domain.User, error)
	GetBalanceView(ctx context.Context, id uuid.UUID) (*domain.BalanceView, error)
	ListSelections(ctx context.Context, date string, strict bool) ([]*domain.Selection, error)
	ListSelectionsFromPicked(
		ctx context.Context,
		stockIDs []string,
		exchangeDate string,
	) ([]*domain.Selection, error)
	GetRealTimeMonitoringKeys(ctx context.Context) ([]string, error)
	GetLatestChip(ctx context.Context) ([]*domain.Selection, error)
	ListOrders(ctx context.Context, arg *domain.ListOrdersParams) ([]*domain.Order, error)
	ListOpenOrders(
		ctx context.Context,
		userID uuid.UUID,
		stockID string,
		orderType string,
	) ([]*domain.Order, error)
	CreateOrder(
		ctx context.Context,
		orders []*domain.Order,
		transactions []*domain.Transaction,
	) error
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) error
	RetrieveDailyCloseHistory(
		ctx context.Context,
		stockIDs []string,
		opts ...string,
	) ([]*domain.DailyClose, error)
	RetrieveThreePrimaryHistory(
		ctx context.Context,
		stockIDs []string,
		opts ...string,
	) ([]*domain.ThreePrimary, error)
	GetHighestPrice(
		ctx context.Context,
		stockIDs []string,
		date string,
		rewindWeek int,
	) (map[string]float32, error)
}

var _ Adapter = (*Imp)(nil)

type Imp struct {
	repo *sqlc.Repo
}

func NewAdapterImp(repo *sqlc.Repo) *Imp {
	return &Imp{
		repo: repo,
	}
}

func (a *Imp) BatchUpsertStocks(ctx context.Context, objs []*domain.Stock) error {
	return a.repo.BatchUpsertStocks(ctx, objs)
}

func (a *Imp) CreateStock(ctx context.Context, obj *domain.Stock) error {
	return a.repo.CreateStock(ctx, obj)
}

func (a *Imp) DeleteStockByID(ctx context.Context, id string) error {
	return a.repo.DeleteStockByID(ctx, id)
}

func (a *Imp) ListCategories(ctx context.Context) ([]*string, error) {
	return a.repo.ListCategories(ctx)
}

func (a *Imp) ListStocks(
	ctx context.Context,
	arg *domain.ListStocksParams,
) ([]*domain.Stock, error) {
	return a.repo.ListStocks(ctx, arg)
}

func (a *Imp) BatchUpsertThreePrimary(ctx context.Context, objs []*domain.ThreePrimary) error {
	return a.repo.BatchUpsertThreePrimary(ctx, objs)
}

func (a *Imp) CreateThreePrimary(ctx context.Context, arg *domain.ThreePrimary) error {
	return a.repo.CreateThreePrimary(ctx, arg)
}

func (a *Imp) ListThreePrimary(
	ctx context.Context,
	arg *domain.ListThreePrimaryParams,
) ([]*domain.ThreePrimary, error) {
	return a.repo.ListThreePrimary(ctx, arg)
}

func (a *Imp) BatchUpsertDailyClose(ctx context.Context, objs []*domain.DailyClose) error {
	return a.repo.BatchUpsertDailyClose(ctx, objs)
}

func (a *Imp) CreateDailyClose(ctx context.Context, obj *domain.DailyClose) error {
	return a.repo.CreateDailyClose(ctx, obj)
}

func (a *Imp) HasDailyClose(ctx context.Context, date string) (bool, error) {
	return a.repo.HasDailyClose(ctx, date)
}

func (a *Imp) ListDailyClose(
	ctx context.Context,
	arg *domain.ListDailyCloseParams,
) ([]*domain.DailyClose, error) {
	return a.repo.ListDailyClose(ctx, arg)
}

func (a *Imp) ListLatestPrice(
	ctx context.Context,
	stockIDs []string,
) ([]*domain.StockPrice, error) {
	return a.repo.ListLatestPrice(ctx, stockIDs)
}

func (a *Imp) BatchUpsertStakeConcentration(
	ctx context.Context,
	objs []*domain.StakeConcentration,
) error {
	return a.repo.BatchUpsertStakeConcentration(ctx, objs)
}

func (a *Imp) GetStakeConcentrationByStockID(
	ctx context.Context,
	stockID,
	date string,
) (*domain.StakeConcentration, error) {
	return a.repo.GetStakeConcentrationByStockID(ctx, stockID, date)
}

func (a *Imp) GetStakeConcentrationsWithVolumes(
	ctx context.Context,
	stockID, date string,
) ([]*domain.CalculationBase, error) {
	return a.repo.GetStakeConcentrationsWithVolumes(ctx, stockID, date)
}

func (a *Imp) HasStakeConcentration(ctx context.Context, exchangeDate string) (bool, error) {
	return a.repo.HasStakeConcentration(ctx, exchangeDate)
}

func (a *Imp) GetStakeConcentrationLatestDataPoint(ctx context.Context) (string, error) {
	return a.repo.GetStakeConcentrationLatestDataPoint(ctx)
}

func (a *Imp) CreatePickedStocks(ctx context.Context, objs []*domain.PickedStock) error {
	return a.repo.CreatePickedStocks(ctx, objs)
}

func (a *Imp) DeletePickedStock(
	ctx context.Context,
	userID uuid.UUID,
	stockID string,
) error {
	return a.repo.DeletePickedStock(ctx, userID, stockID)
}

func (a *Imp) ListPickedStocks(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.PickedStock, error) {
	return a.repo.ListPickedStocks(ctx, userID)
}

func (a *Imp) CreateUser(ctx context.Context, obj *domain.User) error {
	return a.repo.CreateUser(ctx, obj)
}

func (a *Imp) UpdateUser(ctx context.Context, obj *domain.User) error {
	return a.repo.UpdateUser(ctx, obj)
}

func (a *Imp) UpdateSessionID(ctx context.Context, params *domain.UpdateSessionIDParams) error {
	return a.repo.UpdateSessionID(ctx, params)
}

func (a *Imp) DeleteSessionID(ctx context.Context, userID uuid.UUID) error {
	return a.repo.DeleteSessionID(ctx, userID)
}

func (a *Imp) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	return a.repo.DeleteUserByID(ctx, userID)
}

func (a *Imp) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return a.repo.GetUserByEmail(ctx, email)
}

func (a *Imp) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return a.repo.GetUserByPhone(ctx, phone)
}

func (a *Imp) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return a.repo.GetUserByID(ctx, userID)
}

func (a *Imp) ListUsers(ctx context.Context, limit int32, offset int32) ([]*domain.User, error) {
	return a.repo.ListUsers(ctx, limit, offset)
}

func (a *Imp) GetBalanceView(ctx context.Context, id uuid.UUID) (*domain.BalanceView, error) {
	return a.repo.GetBalanceView(ctx, id)
}

func (a *Imp) ListSelections(
	ctx context.Context,
	date string,
	strict bool,
) ([]*domain.Selection, error) {
	return a.repo.ListSelections(ctx, date, strict)
}

func (a *Imp) ListSelectionsFromPicked(
	ctx context.Context,
	stockIDs []string,
	exchangeDate string,
) ([]*domain.Selection, error) {
	return a.repo.ListSelectionsFromPicked(ctx, stockIDs, exchangeDate)
}

func (a *Imp) GetRealTimeMonitoringKeys(ctx context.Context) ([]string, error) {
	return a.repo.GetRealTimeMonitoringKeys(ctx)
}

func (a *Imp) GetLatestChip(ctx context.Context) ([]*domain.Selection, error) {
	return a.repo.GetLatestChip(ctx)
}

func (a *Imp) ListOrders(
	ctx context.Context,
	arg *domain.ListOrdersParams,
) ([]*domain.Order, error) {
	return a.repo.ListOrders(ctx, arg)
}

func (a *Imp) ListOpenOrders(
	ctx context.Context,
	userID uuid.UUID,
	stockID string,
	orderType string,
) ([]*domain.Order, error) {
	return a.repo.ListOpenOrders(ctx, userID, stockID, orderType)
}

func (a *Imp) CreateOrder(
	ctx context.Context,
	orders []*domain.Order,
	transactions []*domain.Transaction,
) error {
	return a.repo.CreateOrder(ctx, orders, transactions)
}

func (a *Imp) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	return a.repo.CreateTransaction(ctx, transaction)
}

func (a *Imp) RetrieveDailyCloseHistory(
	ctx context.Context,
	stockIDs []string,
	opts ...string,
) ([]*domain.DailyClose, error) {
	return a.repo.RetrieveDailyCloseHistory(ctx, stockIDs, opts...)
}

func (a *Imp) RetrieveThreePrimaryHistory(
	ctx context.Context,
	stockIDs []string,
	opts ...string,
) ([]*domain.ThreePrimary, error) {
	return a.repo.RetrieveThreePrimaryHistory(ctx, stockIDs, opts...)
}

func (a *Imp) GetHighestPrice(
	ctx context.Context,
	stockIDs []string,
	date string,
	rewindWeek int,
) (map[string]float32, error) {
	return a.repo.GetHighestPrice(ctx, stockIDs, date, rewindWeek)
}
