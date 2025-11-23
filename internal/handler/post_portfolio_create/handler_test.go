package post_portfolio_create

import (
	"bytes"
	"context"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/yasonofriychuk/tinvest-balancer/internal/generated/api/portfolio_create"
	"github.com/yasonofriychuk/tinvest-balancer/internal/middleware/auth"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest/user_service_client"
	"github.com/yasonofriychuk/tinvest-balancer/internal/testcontainer/postgres"
	"github.com/yasonofriychuk/tinvest-balancer/internal/usecase/portfolio_create"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/crypt"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/logger"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type portfolio struct {
	UserId                 int64      `db:"user_id"`
	AccountId              string     `db:"account_id"`
	Name                   string     `db:"name"`
	AccountType            int64      `db:"account_type"`
	Token                  string     `db:"token"`
	AutoRebalancingEnabled bool       `db:"auto_rebalancing_enabled"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
	DeletedAt              *time.Time `db:"deleted_at"`
}

func TestHandler_PostPortfolioCreate(t *testing.T) {
	var (
		token       = "token"
		accountId   = "09214h"
		userId      = int64(1234)
		currentTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name string

		token     string
		accountId string
		userId    int64

		mockTInvestService func(m *portfolio_create.MocktInvestService)

		statusCode int
		body       any
		portfolio  portfolio
	}{
		{
			name:      "success",
			token:     token,
			accountId: accountId,
			userId:    userId,
			mockTInvestService: func(m *portfolio_create.MocktInvestService) {
				m.EXPECT().GetAccounts(
					gomock.Any(), token, user_service_client.AccountStatus_ACCOUNT_STATUS_OPEN,
				).Return(user_service_client.GetAccountsResponse{
					Accounts: []user_service_client.Account{
						{
							Id:          accountId,
							Name:        "ИИС",
							Type:        user_service_client.AccountType_ACCOUNT_TYPE_TINKOFF_IIS,
							Status:      user_service_client.AccountStatus_ACCOUNT_STATUS_OPEN,
							AccessLevel: user_service_client.AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS,
						},
					},
				}, nil)

				m.EXPECT().MarginIsEnabled(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			statusCode: http.StatusOK,
			body: api.PostPortfolioCreate200JSONResponse{
				Status: "Успех",
			},
			portfolio: portfolio{
				UserId:                 userId,
				AccountId:              accountId,
				Name:                   "ИИС",
				AccountType:            int64(user_service_client.AccountType_ACCOUNT_TYPE_TINKOFF_IIS),
				Token:                  token,
				AutoRebalancingEnabled: false,
				CreatedAt:              currentTime,
				UpdatedAt:              currentTime,
				DeletedAt:              nil,
			},
		},
	}

	ctx := context.Background()
	ctr := postgres.MustNew(ctx)

	defer func() {
		require.NoError(t, ctr.Terminate())
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				require.NoError(t, ctr.Restore(ctx))
			})

			ctx := context.WithoutCancel(ctx)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, err := ctr.GetDb(ctx)
			require.NoError(t, err)

			cryptKey, err := crypt.GenerateBase64Key()
			require.NoError(t, err)
			crypter := crypt.MustNew(crypt.MustKeyFromBase64(cryptKey))

			defer func(db *sqlx.DB) {
				require.NoError(t, db.Close())
			}(db)

			tinvestService := portfolio_create.NewMocktInvestService(ctrl)
			tt.mockTInvestService(tinvestService)

			portfolioCreateUsecase := portfolio_create.New(
				portfolio_create.NewStorage(db).WithNowFunc(func() time.Time { return currentTime }),
				crypter,
				tinvestService,
			)

			h := New(logger.NewLogger(slog.LevelDebug, "test", io.Discard), portfolioCreateUsecase)
			ts := httptest.NewServer(
				auth.NewAuthMiddleware().Handler(
					api.Handler(api.NewStrictHandler(h, nil)),
				),
			)
			defer ts.Close()

			inputData := api.PostPortfolioCreateJSONBody{
				Token: tt.token,
			}

			inputJson, err := json.Marshal(inputData)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/portfolio/create", bytes.NewBuffer(inputJson))
			require.NoError(t, err)

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", strconv.FormatInt(tt.userId, 10))

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			defer func(Body io.ReadCloser) {
				require.NoError(t, Body.Close())
			}(resp.Body)

			expected, err := json.Marshal(tt.body)
			require.NoError(t, err)

			actual, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			assert.JSONEq(t, string(expected), string(actual))

			var portfolios []portfolio
			require.NoError(t, db.Select(
				&portfolios,
				`
				SELECT 
				    user_id, 
				    account_id, 
				    name, 
				    account_type, 
				    token, 
				    auto_rebalancing_enabled, 
				    created_at, 
				    updated_at, 
				    deleted_at 
				FROM portfolio WHERE account_id = $1`,
				tt.accountId,
			))

			require.Len(t, portfolios, 1)

			tokenDecrypt, err := crypter.Decrypt(portfolios[0].Token)
			require.NoError(t, err)
			portfolios[0].Token = tokenDecrypt

			assert.Equal(t, tt.portfolio, portfolios[0])
		})
	}
}
