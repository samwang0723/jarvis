package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"golang.org/x/crypto/bcrypt"
)

func (s *serviceImpl) Login(ctx context.Context, email, password string) (obj *entity.User, err error) {
	obj, err = s.dal.GetUserByEmail(ctx, email)
	if err != nil || obj == nil {
		return nil, errUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(password))
	if err != nil {
		return nil, errUserPasswordNotMatch
	}

	// generate session_id
	err = s.dal.UpdateSessionID(ctx, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
