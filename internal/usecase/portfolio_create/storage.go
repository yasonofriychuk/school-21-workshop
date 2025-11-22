package portfolio_create

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db  *sqlx.DB
	now func() time.Time
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (s *Storage) setToken(cxt context.Context, userId int64, accountId string, name string, accountType int64, token string) error {
	const query = `
		INSERT INTO portfolio 
		    (user_id, account_id, name, account_type, token) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT 
		    (user_id, account_id) 
		DO UPDATE 
		    SET token = EXCLUDED.token, 
		        name = EXCLUDED.name,
		        account_type = EXCLUDED.account_type,
		        updated_at = DEFAULT,
		        deleted_at = NULL,
		        auto_rebalancing_enabled = false
	`

	if _, err := s.db.ExecContext(cxt, query, userId, accountId, name, accountType, token); err != nil {
		return fmt.Errorf("s.db.ExecContext: %w", err)
	}
	return nil
}
