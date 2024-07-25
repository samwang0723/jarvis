package middleware_test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/golang/mock/gomock"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/middleware"
	mock_services "github.com/samwang0723/jarvis/internal/app/services/mocks"
	"github.com/samwang0723/jarvis/internal/common/testhelper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	testhelper.LoadTestConfig()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func createTestToken(userID, sessionID string) string {
	signer, _ := jwt.NewSignerHS(jwt.HS256, []byte(config.GetCurrentConfig().JwtSecret))
	builder := jwt.NewBuilder(signer)

	claims := jwt.RegisteredClaims{
		ID:        sessionID,
		Subject:   userID,
		Audience:  []string{"jarvis"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token, _ := builder.Build(claims)
	return token.String()
}

func TestAuthenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockIService(ctrl)
	authFunc := middleware.Authenticate(mockService)
	// Create a valid token
	userID := uuid.Must(uuid.NewV4())
	sessionID := uuid.Must(uuid.NewV4())

	type args struct {
		token     string
		userID    uuid.UUID
		sessionID uuid.UUID
	}

	tests := []struct {
		name              string
		args              args
		wantErr           bool
		expectServiceCall bool
	}{
		{
			name: "parse token successfully",
			args: args{
				userID:    userID,
				sessionID: sessionID,
				token:     createTestToken(userID.String(), sessionID.String()),
			},
			wantErr:           false,
			expectServiceCall: true,
		},
		{
			name:              "parse token failed",
			args:              args{userID: uuid.Nil, sessionID: uuid.Nil, token: "invalid_token"},
			wantErr:           true,
			expectServiceCall: false,
		},
		{
			name: "parse subject failed",
			args: args{
				userID:    uuid.Nil,
				sessionID: uuid.Nil,
				token:     createTestToken("failed_userID", sessionID.String()),
			},
			wantErr:           true,
			expectServiceCall: false,
		},
		{
			name: "parse sessionID failed",
			args: args{
				userID:    uuid.Nil,
				sessionID: uuid.Nil,
				token:     createTestToken(userID.String(), "failed_sessionID"),
			},
			wantErr:           true,
			expectServiceCall: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.expectServiceCall {
				// Set up mock expectations
				mockService.EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(&domain.User{
						ID:        domain.ID{ID: tt.args.userID},
						SessionID: tt.args.sessionID.String(),
					}, nil).Times(1)
			}

			// Create a context with the token
			ctx := metadata.NewIncomingContext(
				context.Background(),
				metadata.Pairs("authorization", "bearer "+tt.args.token),
			)

			// Call the function
			newCtx, err := authFunc(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, newCtx)
			} else {
				assert.NotNil(t, newCtx)
				assert.NoError(t, err)
				user := newCtx.Value(config.JwtClaimsKey).(*domain.User)
				assert.Equal(t, userID, user.ID.ID)
				assert.Equal(t, sessionID.String(), user.SessionID)
			}
		})
	}
}
