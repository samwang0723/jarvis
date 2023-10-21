package handlers

import (
	"context"

	"github.com/cristalhq/jwt/v5"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/helper"
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

	// create claims (you can create your own, see: ExampleBuilder_withUserClaims)
	claims := &jwt.RegisteredClaims{
		Audience:  []string{"jarvis"},
		ID:        user.SessionID,
		ExpiresAt: jwt.NewNumericDate(*user.SessionExpiredAt),
		Issuer:    user.Email,
		Subject:   helper.Uint64ToString(user.ID.Uint64()),
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
