package handlers

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/helper"
)

func (h *handlerImpl) Login(ctx context.Context, req *dto.LoginRequest) *dto.LoginResponse {
	user, err := h.dataService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &dto.LoginResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}
	}

	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: user.SessionExpiredAt.Unix(),
		Issuer:    user.Email,
		Id:        user.SessionID,
		Subject:   helper.Uint64ToString(user.ID.Uint64()),
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(h.jwtSecret)
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
		AccessToken:  tokenString,
	}
}
