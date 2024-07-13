package sqlc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

func (repo *Repo) CreateUser(ctx context.Context, obj *domain.User) error {
	return repo.RunInTransaction(ctx, func(ctx context.Context) error {
		obj.ID.ID = uuid.Must(uuid.NewV4())
		err := repo.primary().CreateUser(ctx, &sqlcdb.CreateUserParams{
			ID:        obj.ID.ID,
			FirstName: obj.FirstName,
			LastName:  obj.LastName,
			Email:     obj.Email,
			Phone:     obj.Phone,
			Password:  obj.Password,
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		err = repo.createBalance(ctx, obj.ID.ID, 0.0)
		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *Repo) UpdateUser(ctx context.Context, obj *domain.User) error {
	return repo.primary().UpdateUser(ctx, &sqlcdb.UpdateUserParams{
		ID:        obj.ID.ID,
		FirstName: obj.FirstName,
		LastName:  obj.LastName,
		Email:     obj.Email,
		Phone:     obj.Phone,
		Password:  obj.Password,
	})
}

func (repo *Repo) UpdateSessionID(ctx context.Context, params *domain.UpdateSessionIDParams) error {
	return repo.primary().UpdateSessionID(ctx, &sqlcdb.UpdateSessionIDParams{
		SessionID:        sql.NullString{String: params.SessionID, Valid: true},
		SessionExpiredAt: sql.NullTime{Time: params.SessionExpiredAt, Valid: true},
		ID:               params.ID,
	})
}

func (repo *Repo) DeleteSessionID(ctx context.Context, userID uuid.UUID) error {
	return repo.primary().DeleteSessionID(ctx, userID)
}

func (repo *Repo) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	return repo.primary().DeleteUserByID(ctx, userID)
}

func (repo *Repo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := repo.primary().GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDomainUser(row), nil
}

func (repo *Repo) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	row, err := repo.primary().GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return toDomainUser(row), nil
}

func (repo *Repo) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	row, err := repo.primary().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toDomainUser(row), nil
}

func (repo *Repo) ListUsers(
	ctx context.Context,
	limit int32,
	offset int32,
) ([]*domain.User, error) {
	result, err := repo.primary().ListUsers(ctx, &sqlcdb.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUserList(result), nil
}

func toDomainUserList(rows []*sqlcdb.User) []*domain.User {
	result := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomainUser(row))
	}
	return result
}

func toDomainUser(row *sqlcdb.User) *domain.User {
	return &domain.User{
		ID:               domain.ID{ID: row.ID},
		FirstName:        row.FirstName,
		LastName:         row.LastName,
		Email:            row.Email,
		Phone:            row.Phone,
		Password:         row.Password,
		SessionExpiredAt: &row.SessionExpiredAt.Time,
		PhoneConfirmedAt: &row.PhoneConfirmedAt.Time,
		EmailConfirmedAt: &row.EmailConfirmedAt.Time,
		Time: domain.Time{
			CreatedAt: &row.CreatedAt,
			UpdatedAt: &row.UpdatedAt,
		},
	}
}
