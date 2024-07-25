package services

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionExpiredDays = 5
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
	newSessionID := uuid.Must(uuid.NewV4())
	newExpiredAt := time.Now().AddDate(0, 0, sessionExpiredDays)
	err = s.dal.UpdateSessionID(ctx, &domain.UpdateSessionIDParams{
		ID:               obj.ID.ID,
		SessionID:        newSessionID.String(),
		SessionExpiredAt: newExpiredAt,
	})
	obj.SessionID = newSessionID.String()
	obj.SessionExpiredAt = &newExpiredAt
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *serviceImpl) getCurrentUserID(ctx context.Context) (userID uuid.UUID, err error) {
	user, ok := ctx.Value(config.JwtClaimsKey).(*domain.User)
	if !ok {
		return uuid.Nil, errUserNotFound
	}

	return user.ID.ID, nil
}

func (s *serviceImpl) Logout(ctx context.Context) error {
	return s.dal.DeleteSessionID(ctx, s.currentUserID)
}
