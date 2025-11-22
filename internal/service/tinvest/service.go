package tinvest

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"github.com/russianinvestments/invest-api-go-sdk/retry"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest/user_service_client"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	waitBetween   = 500 * time.Millisecond
	headerAppName = "x-app-name"
)

type Service struct {
	log  logger.Log
	conn *grpc.ClientConn
	*user_service_client.UsersServiceClient
}

func MustNew(log logger.Log, appName string) *Service {
	conn, err := newConn(investgo.Config{
		EndPoint:                      "invest-public-api.tinkoff.ru:443",
		AppName:                       appName,
		DisableResourceExhaustedRetry: false,
		DisableAllRetry:               false,
		MaxRetries:                    3,
	})
	if err != nil {
		panic(fmt.Errorf("create tinvest service failed: %v", err))
	}

	return &Service{
		log:                log,
		conn:               conn,
		UsersServiceClient: user_service_client.New(conn),
	}
}

func (s *Service) Stop() error {
	return s.conn.Close()
}

func newConn(conf investgo.Config) (*grpc.ClientConn, error) {
	var dialOpts []grpc.DialOption

	opts := []retry.CallOption{
		retry.WithCodes(codes.Unavailable, codes.Internal, codes.Canceled),
		retry.WithBackoff(retry.BackoffLinear(waitBetween)),
		retry.WithMax(3),
	}

	// при исчерпывании лимита запросов в минуту, нужно ждать дольше
	exhaustedOpts := []retry.CallOption{
		retry.WithCodes(codes.ResourceExhausted),
		retry.WithMax(conf.MaxRetries),
		retry.WithOnRetryCallback(func(ctx context.Context, attempt uint, err error) {
			// TODO тут можно добавить лог
		}),
	}

	streamInterceptors := []grpc.StreamClientInterceptor{
		retry.StreamClientInterceptor(opts...),
		outgoingAppNameStreamInterceptor(conf.AppName),
	}

	var unaryInterceptors []grpc.UnaryClientInterceptor
	if conf.DisableResourceExhaustedRetry {
		unaryInterceptors = []grpc.UnaryClientInterceptor{
			retry.UnaryClientInterceptor(opts...),
			outgoingAppNameUnaryInterceptor(conf.AppName),
		}
	} else {
		unaryInterceptors = []grpc.UnaryClientInterceptor{
			retry.UnaryClientInterceptor(opts...),
			retry.UnaryClientInterceptorRE(exhaustedOpts...),
			outgoingAppNameUnaryInterceptor(conf.AppName),
		}
	}

	dialOpts = append(
		dialOpts,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), // nolint: gosec
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
	)

	conn, err := grpc.NewClient(conf.EndPoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return conn, nil
}

func outgoingAppNameUnaryInterceptor(appName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, headerAppName, appName)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func outgoingAppNameStreamInterceptor(appName string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, headerAppName, appName)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
