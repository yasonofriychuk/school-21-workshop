//go:generate go tool mockgen -source=$GOFILE -destination contract_mock.go -package $GOPACKAGE

package portfolio_create

import (
	"context"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest/user_service_client"
)

type storage interface {
	setToken(cxt context.Context, userId int64, accountId string, name string, accountType int64, token string) error
}

type crypter interface {
	Encrypt(text string) (string, error)
}

type tInvestService interface {
	GetAccounts(ctx context.Context, token string, status user_service_client.AccountStatus) (user_service_client.GetAccountsResponse, error)
	MarginIsEnabled(ctx context.Context, token string, accountId string) (bool, error)
}
