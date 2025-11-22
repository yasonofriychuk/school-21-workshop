package portfolio_create

import (
	"context"
	"errors"
	"fmt"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest/user_service_client"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/cerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"slices"
)

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
	resp, err := u.tInvestService.GetAccounts(ctx, token, user_service_client.AccountStatus_ACCOUNT_STATUS_OPEN)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			return cerrors.NewUserError(err).WithMessage("Токен недействителен")
		}
		return fmt.Errorf("u.tInvestService.GetAccounts: %w", err)
	}

	if len(resp.Accounts) == 0 {
		return cerrors.NewUserError(errors.New("accounts is empty")).WithMessage("Нет доступных счетов")
	}

	if len(resp.Accounts) > 1 {
		return cerrors.NewUserError(errors.New("accounts more than 1")).WithMessage("Токен должен быть создан для одного счета")
	}

	acc := resp.Accounts[0]
	if acc.AccessLevel != user_service_client.AccessLevel_ACCOUNT_ACCESS_LEVEL_FULL_ACCESS {
		return cerrors.NewUserError(errors.New("there are no accounts with full access")).WithMessage("Для совершения операций нужен токен с полным доступом")
	}

	if !slices.Contains([]user_service_client.AccountType{
		user_service_client.AccountType_ACCOUNT_TYPE_TINKOFF, user_service_client.AccountType_ACCOUNT_TYPE_TINKOFF_IIS,
	}, acc.Type) {
		return cerrors.NewUserError(errors.New("incorrect account type")).WithMessage("Выставление ордеров для этого счета недоступно")
	}

	if acc.Status != user_service_client.AccountStatus_ACCOUNT_STATUS_OPEN {
		return cerrors.NewUserError(errors.New("incorrect account status")).WithMessage("Счет закрыт")
	}

	enabled, err := u.tInvestService.MarginIsEnabled(ctx, token, acc.Id)
	if err != nil {
		return fmt.Errorf("u.tInvestService.GetMarginAttributes: %w", err)
	}

	if enabled {
		return cerrors.NewUserError(errors.New("margin enabled")).WithMessage("Отключите маржинальную торговлю")
	}

	// Шифрование и сохранение токена
	encryptedToken, err := u.crypter.Encrypt(token)
	if err != nil {
		return fmt.Errorf("u.crypter.Encrypt: %w", err)
	}

	if err := u.s.setToken(ctx, userId, acc.Id, acc.Name, convertAccountType(acc.Type), encryptedToken); err != nil {
		return fmt.Errorf("u.s.setToken: %w", err)
	}

	return nil
}

func convertAccountType(tType user_service_client.AccountType) int64 {
	switch tType {
	case user_service_client.AccountType_ACCOUNT_TYPE_TINKOFF:
		return 1
	case user_service_client.AccountType_ACCOUNT_TYPE_TINKOFF_IIS:
		return 2
	case user_service_client.AccountType_ACCOUNT_TYPE_INVEST_BOX:
		return 3
	case user_service_client.AccountType_ACCOUNT_TYPE_INVEST_FUND:
		return 4
	default:
		return 0
	}
}
