package middleware

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/pb"
	"github.com/samwang0723/jarvis/internal/app/services"
	"github.com/samwang0723/jarvis/internal/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseToken(token string) (*jwt.RegisteredClaims, error) {
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(config.GetCurrentConfig().JwtSecret))
	if err != nil {
		return nil, err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return nil, err
	}

	// or just verify it's signature
	err = verifier.Verify(newToken)
	if err != nil {
		return nil, err
	}

	// get Registered claims
	newClaims := &jwt.RegisteredClaims{}
	err = json.Unmarshal(newToken.Claims(), newClaims)
	if err != nil {
		return nil, err
	}

	err = jwt.ParseClaims(tokenBytes, verifier, newClaims)
	if err != nil {
		return nil, err
	}

	return newClaims, nil
}

// Authenticate is used by a middleware to authenticate requests
func Authenticate(service services.IService) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		tokenInfo, err := parseToken(token)
		if err != nil || !tokenInfo.IsValidAt(time.Now()) || !tokenInfo.IsForAudience("jarvis") {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid auth token: %v", err)
		}

		sessionID := tokenInfo.ID
		userID, err := uuid.FromString(tokenInfo.Subject)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid auth token: %v", err)
		}

		user, err := service.GetUserByID(ctx, userID)
		if err != nil || user.SessionID != sessionID {
			return nil, status.Error(
				codes.Unauthenticated,
				"Invalid auth token: session_id invalid",
			)
		}

		ctx = logging.InjectFields(ctx, logging.Fields{"auth.sub", tokenInfo.Subject})

		return context.WithValue(ctx, config.JwtClaimsKey, user), nil
	}
}

func AuthRoutes(_ context.Context, callMeta interceptors.CallMeta) bool {
	if helper.StringInSlice(callMeta.Method, []string{"Login", "CreateUser"}) {
		return false
	}

	return pb.JarvisV1_ServiceDesc.ServiceName == callMeta.Service
}
