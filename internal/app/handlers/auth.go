package handlers

import (
	"context"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/samwang0723/jarvis/internal/app/dto"
)

func (h *handlerImpl) Login(ctx context.Context, req *dto.LoginRequest) *dto.LoginResponse {
	user, err := h.dataService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &dto.LoginResponse{
			Status:       dto.StatusUnauthorized,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}
	}

	signer, err := jwt.NewSignerHS(jwt.HS256, h.jwtSecret)
	if err != nil {
		return &dto.LoginResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}
	}

	// create claims
	claims := &jwt.RegisteredClaims{
		Audience:  []string{"jarvis"},
		ID:        user.SessionID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(*user.SessionExpiredAt),
		Issuer:    user.Email,
		Subject:   user.ID.ID.String(),
	}

	// create a Builder
	builder := jwt.NewBuilder(signer)

	// and build a Token
	token, err := builder.Build(claims)
	if err != nil {
		return &dto.LoginResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}
	}

	return &dto.LoginResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
		AccessToken:  token.String(),
	}
}

func (h *handlerImpl) Logout(ctx context.Context) *dto.LogoutResponse {
	err := h.dataService.WithUserID(ctx).Logout(ctx)
	if err != nil {
		return &dto.LogoutResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}
	}

	return &dto.LogoutResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}
}
