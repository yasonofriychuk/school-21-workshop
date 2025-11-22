package convertor

import (
	"fmt"

	"github.com/yasonofriychuk/tinvest-balancer/pkg/cerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertErr(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return cerrors.NewInternalError(fmt.Errorf("unknown non-grpc error: %w", err))
	}

	switch st.Code() {
	case codes.Unauthenticated:
		return cerrors.NewBusinessError(fmt.Errorf("invalid token: %s", st.Message()))
	default:
		return cerrors.NewInternalError(st.Err())
	}
}
