package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/rs/cors"
	"github.com/yasonofriychuk/tinvest-balancer/internal/config"
	"github.com/yasonofriychuk/tinvest-balancer/internal/generated/api"
	"github.com/yasonofriychuk/tinvest-balancer/internal/handler/post_portfolio_create"
	"github.com/yasonofriychuk/tinvest-balancer/internal/middleware/auth"
	"github.com/yasonofriychuk/tinvest-balancer/internal/service/tinvest"
	"github.com/yasonofriychuk/tinvest-balancer/internal/usecase/portfolio_create"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/crypt"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/logger"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/migrate"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/postgres"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Run(ctx context.Context) error {
	cfg := config.MustNewConfig()

	log := logger.NewLogger(slog.LevelDebug, cfg.Env, os.Stdout)
	db := postgres.MustNew(cfg.DSN)

	if err := migrate.Migrate(cfg.DSN); err != nil {
		return fmt.Errorf("migrate.Migrate: %w", err)
	}

	crypter := crypt.MustNew(crypt.MustKeyFromBase64(cfg.CryptKey))
	authMW := auth.NewAuthMiddleware()
	tinvestService := tinvest.MustNew(log, "@tinvest_balancer_bot")
	defer func(tinvestService *tinvest.Service) { _ = tinvestService.Stop() }(tinvestService)

	portfolioCreateUsecase := portfolio_create.New(portfolio_create.NewStorage(db), crypter, tinvestService)

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Use(authMW.Handler)
		api.Routers{
			PortfolioCreate: post_portfolio_create.New(log, portfolioCreateUsecase),
		}.Register(r, []nethttp.StrictHTTPMiddlewareFunc{})
	})

	s := http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           cors.New(cors.Options{AllowedOrigins: []string{}}).Handler(router),
		ReadHeaderTimeout: time.Minute * 5,
	}

	errorGroup, _ := errgroup.WithContext(ctx)
	errorGroup.Go(func() (err error) {
		<-ctx.Done()
		log.WithContext(ctx).Info("shutting down server")
		if err := s.Shutdown(ctx); err != nil {
			return fmt.Errorf("s.Shutdown(): %w", err)
		}
		return nil
	})

	errorGroup.Go(func() (err error) {
		log.WithContext(ctx).WithFields(map[string]any{
			"addr": s.Addr,
		}).Info("starting server")

		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("s.ListenAndServe(): %w", err)
		}
		return nil
	})

	if err := errorGroup.Wait(); err != nil {
		return fmt.Errorf("errorGroup.Wait(): %w", err)
	}

	return nil
}
