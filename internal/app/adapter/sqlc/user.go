package sqlc

import (
	"context"
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/domain"
	sqlcdb "github.com/samwang0723/jarvis/internal/db/main/sqlc"
)

func (repo *Repo) CreateUser(ctx context.Context, obj *domain.User) error {
	return repo.RunInTransaction(ctx, func(ctx context.Context) error {
		err := repo.primaryConn.queries.CreateUser(ctx, &sqlcdb.CreateUserParams{
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

		err = repo.CreateBalanceView(ctx, obj.ID.ID, 0.0)
		if err != nil {
			return err
		}

		return nil
	})
}
