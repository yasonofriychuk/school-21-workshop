package user_service_client

import (
	"context"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
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
func (us *UsersServiceClient) GetAccounts(ctx context.Context, token string, status *pb.AccountStatus) (*investgo.GetAccountsResponse, error) {
	var header, trailer metadata.MD
	resp, err := us.pbClient.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{
			Status: status,
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

	return &investgo.GetAccountsResponse{
		GetAccountsResponse: resp,
		Header:              header,
	}, err
}

// GetMarginAttributes - Расчёт маржинальных показателей по счёту
func (us *UsersServiceClient) GetMarginAttributes(ctx context.Context, token string, accountId string) (*investgo.GetMarginAttributesResponse, error) {
	var header, trailer metadata.MD
	resp, err := us.pbClient.GetMarginAttributes(ctx, &pb.GetMarginAttributesRequest{
		AccountId: accountId,
	}, grpc.PerRPCCredentials(oauth.TokenSource{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	}), grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &investgo.GetMarginAttributesResponse{
		GetMarginAttributesResponse: resp,
		Header:                      header,
	}, err
}
