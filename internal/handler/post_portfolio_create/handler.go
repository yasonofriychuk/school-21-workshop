package post_portfolio_create

import (
	"context"

	"github.com/yasonofriychuk/tinvest-balancer/pkg/logger"

	"github.com/AlekSi/pointer"

	api "github.com/yasonofriychuk/tinvest-balancer/internal/generated/api/portfolio_create"
	"github.com/yasonofriychuk/tinvest-balancer/internal/middleware/auth"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/cerrors"
)

type Handler struct {
	log                    logger.Log
	createPortfolioUsecase createPortfolioUsecase
}

func New(log logger.Log, createPortfolioUsecase createPortfolioUsecase) *Handler {
	return &Handler{
		log:                    log,
		createPortfolioUsecase: createPortfolioUsecase,
	}
}

func (h *Handler) PostPortfolioCreate(ctx context.Context, request api.PostPortfolioCreateRequestObject) (api.PostPortfolioCreateResponseObject, error) {
	userId, _ := auth.GetUserId(ctx)
	body := pointer.Get(request.Body)

	if err := h.createPortfolioUsecase.CreatePortfolio(
		ctx,
		userId,
		body.Token,
	); err != nil {
		if cErr, ok := cerrors.AsUserError(err); ok {
			return api.PostPortfolioCreate400JSONResponse{Message: cErr.Message()}, nil
		}
		h.log.WithContext(ctx).WithError(err).WithFields(map[string]any{
			"userId": userId,
		}).Error("failed to create portfolio")
		return api.PostPortfolioCreate500JSONResponse{Message: "Не удалось сохранить токен"}, nil
	}

	return api.PostPortfolioCreate200JSONResponse{Status: "Успех"}, nil
}
