package post_portfolio_create

import "context"

type createPortfolioUsecase interface {
	CreatePortfolio(ctx context.Context, userId int64, token string) error
}
