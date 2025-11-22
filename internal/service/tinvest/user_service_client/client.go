package user_service_client

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest/convertor"
	"google.golang.org/grpc/status"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

type UsersServiceClient struct {
	pbClient pb.UsersServiceClient
}

func New(conn *grpc.ClientConn) *UsersServiceClient {
	return &UsersServiceClient{
		pbClient: pb.NewUsersServiceClient(conn),
	}
}

// GetAccounts - Метод получения счетов пользователя
func (us *UsersServiceClient) GetAccounts(ctx context.Context, token string, status AccountStatus) (GetAccountsResponse, error) {
	var header, trailer metadata.MD
	resp, err := us.pbClient.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{
			Status: pointer.To(pb.AccountStatus(status)),
		},
		grpc.PerRPCCredentials(oauth.TokenSource{
			TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		}),
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		header = trailer
	}

	if resp == nil {
		return GetAccountsResponse{}, convertor.ConvertErr(err)
	}

	return GetAccountsResponse{
		Accounts: lo.FilterMap(resp.GetAccounts(), func(a *pb.Account, _ int) (Account, bool) {
			if a == nil {
				return Account{}, false
			}

			return Account{
				Id:          a.GetId(),
				Name:        a.GetName(),
				Type:        AccountType(a.GetType()),
				AccessLevel: AccessLevel(a.GetAccessLevel()),
				Status:      AccountStatus(a.GetStatus()),
			}, true
		}),
	}, err
}

// MarginIsEnabled - Проверка наличия включенной маржинальной торговли
func (us *UsersServiceClient) MarginIsEnabled(ctx context.Context, token string, accountId string) (bool, error) {
	_, err := us.pbClient.GetMarginAttributes(ctx, &pb.GetMarginAttributesRequest{
		AccountId: accountId,
	}, grpc.PerRPCCredentials(oauth.TokenSource{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	}))

	if s, ok := status.FromError(err); ok && s != nil {
		if s.Message() == "30051" {
			return false, nil
		}
	}

	if err != nil {
		return false, convertor.ConvertErr(err)
	}

	return true, nil
}
