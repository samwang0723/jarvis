package sqlc

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
)

func (repo *Repo) CreateUser(ctx context.Context, obj *domain.User) (err error) {
	return nil
}
