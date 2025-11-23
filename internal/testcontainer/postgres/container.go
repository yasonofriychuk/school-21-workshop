package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/yasonofriychuk/tinvest-balancer/pkg/migrate"
	postgres_db "github.com/yasonofriychuk/tinvest-balancer/pkg/postgres"
)

type Container struct {
	*postgres.PostgresContainer
}

func MustNew(ctx context.Context) Container {
	ctr, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("master"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		panic(fmt.Errorf("postgres.Run: %w", err))
	}

	dsnUrl, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("ctr.ConnectionString: %w", err))
	}

	if err := migrate.Migrate(dsnUrl, "file://../../../migrations"); err != nil {
		panic(fmt.Errorf("migrate.Migrate: %w", err))
	}

	if err := ctr.Snapshot(ctx); err != nil {
		panic(fmt.Errorf("ctr.Snapshot(ctx): %w", err))
	}

	return Container{
		PostgresContainer: ctr,
	}
}

func (c *Container) GetDb(ctx context.Context) (*sqlx.DB, error) {
	dsnUrl, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("c.ConnectionString: %w", err)
	}

	db, err := postgres_db.New(dsnUrl)
	if err != nil {
		return nil, fmt.Errorf("postgres_db.New: %w", err)
	}
	return db, nil
}

func (c *Container) Terminate() error {
	return testcontainers.TerminateContainer(c.Container)
}
