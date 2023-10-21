package services

import (
	"context"

	"github.com/cristalhq/jwt/v5"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/app/middleware"
	"github.com/samwang0723/jarvis/internal/helper"
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

func (s *serviceImpl) getCurrentUserID(ctx context.Context) (userID uint64, err error) {
	// get user_id from context
	claims, ok := ctx.Value(middleware.JwtClaimsKey).(jwt.RegisteredClaims)
	if !ok {
		return 0, errInvalidJWTToken
	}

	sessionID := claims.ID
	userID, err = helper.StringToUint64(claims.Subject)
	if err != nil {
		return 0, err
	}

	user, err := s.GetUser(ctx)
	if err != nil || user.SessionID != sessionID {
		return 0, err
	}

	return userID, nil
}

func (s *serviceImpl) Logout(ctx context.Context) error {
	return s.dal.DeleteSessionID(ctx, s.currentUserID)
}
