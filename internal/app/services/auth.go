package services

import (
	"context"

	"github.com/cristalhq/jwt/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/middleware"
	"golang.org/x/crypto/bcrypt"
)

func (s *serviceImpl) Login(
	ctx context.Context,
	email, password string,
) (obj *domain.User, err error) {
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

func (s *serviceImpl) getCurrentUserID(ctx context.Context) (userID uuid.UUID, err error) {
	// get user_id from context
	claims, ok := ctx.Value(middleware.JwtClaimsKey).(jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, errInvalidJWTToken
	}

	s.logger.Info().Msgf("claims: %+v", claims)

	sessionID := claims.ID
	userID, err = uuid.FromString(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	user, err := s.dal.GetUserByID(ctx, userID)
	if err != nil || user.SessionID != sessionID {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *serviceImpl) Logout(ctx context.Context) error {
	return s.dal.DeleteSessionID(ctx, s.currentUserID)
}
