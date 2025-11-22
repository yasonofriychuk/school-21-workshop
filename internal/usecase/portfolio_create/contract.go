package portfolio_create

import (
	"context"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type storage interface {
	setToken(cxt context.Context, userId int64, accountId string, name string, accountType int64, token string) error
}

type crypter interface {
	Encrypt(text string) (string, error)
}

type tInvestService interface {
	GetAccounts(ctx context.Context, token string, status *pb.AccountStatus) (*investgo.GetAccountsResponse, error)
	GetMarginAttributes(ctx context.Context, token string, accountId string) (*investgo.GetMarginAttributesResponse, error)
}
