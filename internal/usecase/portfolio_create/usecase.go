package portfolio_create

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/AlekSi/pointer"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/cerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const accountMarginDisabledCode = "30051"

type Usecase struct {
	s              storage
	crypter        crypter
	tInvestService tInvestService
}

func New(storage storage, crypter crypter, tInvestService tInvestService) *Usecase {
	return &Usecase{
		s:              storage,
		crypter:        crypter,
		tInvestService: tInvestService,
	}
}

func (u *Usecase) CreatePortfolio(ctx context.Context, userId int64, token string) error {
	// Проверка токена
	resp, err := u.tInvestService.GetAccounts(ctx, token, pointer.To(pb.AccountStatus_ACCOUNT_STATUS_OPEN))
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			return cerrors.NewUserError(err).WithMessage("Токен недействителен")
		}
		return tinvest.ConvertErr(fmt.Errorf("u.tInvestService.GetAccounts"))
	}

	if len(resp.GetAccounts()) == 0 {
		return cerrors.NewUserError(errors.New("accounts is empty")).WithMessage("Нет доступных счетов")
	}

	if len(resp.GetAccounts()) > 1 {
		return cerrors.NewUserError(errors.New("accounts more than 1")).WithMessage("Токен должен быть создан для одного счета")
	}

	acc := resp.GetAccounts()[0]
	if acc == nil {
		return cerrors.NewUserError(errors.New("accounts is empty")).WithMessage("Нет доступных счетов")
	}

	if acc.GetAccessLevel() != pb.AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS {
		return cerrors.NewUserError(errors.New("there are no accounts with full access")).WithMessage("Для совершения операций нужен токен с полным доступом")
	}

	if !slices.Contains([]pb.AccountType{
		pb.AccountType_ACCOUNT_TYPE_TINKOFF, pb.AccountType_ACCOUNT_TYPE_TINKOFF_IIS,
	}, acc.GetType()) {
		return cerrors.NewUserError(errors.New("incorrect account type")).WithMessage("Выставление ордеров для этого счета недоступно")
	}

	if acc.GetStatus() != pb.AccountStatus_ACCOUNT_STATUS_OPEN {
		return cerrors.NewUserError(errors.New("incorrect account status")).WithMessage("Счет закрыт")
	}

	_, err = u.tInvestService.GetMarginAttributes(ctx, token, acc.GetId())
	if s, ok := status.FromError(err); ok {
		if s.Message() != accountMarginDisabledCode {
			return cerrors.NewUserError(errors.New("margin enabled")).WithMessage("Отключите маржинальную торговлю")
		}
	} else {
		return fmt.Errorf("u.tInvestService.GetMarginAttributes: %w", err)
	}

	// Шифрование и сохранение токена
	encryptedToken, err := u.crypter.Encrypt(token)
	if err != nil {
		return fmt.Errorf("u.crypter.Encrypt: %w", err)
	}

	if err := u.s.setToken(ctx, userId, acc.GetId(), acc.GetName(), convertAccountType(acc.GetType()), encryptedToken); err != nil {
		return fmt.Errorf("u.s.setToken: %w", err)
	}

	return nil
}

func convertAccountType(tType pb.AccountType) int64 {
	switch tType {
	case pb.AccountType_ACCOUNT_TYPE_TINKOFF:
		return 1
	case pb.AccountType_ACCOUNT_TYPE_TINKOFF_IIS:
		return 2
	case pb.AccountType_ACCOUNT_TYPE_INVEST_BOX:
		return 3
	case pb.AccountType_ACCOUNT_TYPE_INVEST_FUND:
		return 4
	default:
		return 0
	}
}
