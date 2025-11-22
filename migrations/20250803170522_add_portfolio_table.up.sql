CREATE TABLE IF NOT EXISTS portfolio (
    user_id BIGINT NOT NULL,
    account_id TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    account_type smallint not null default 1,
    token TEXT NOT NULL,
    auto_rebalancing_enabled BOOLEAN default false,

    created_at timestamp default current_timestamp NOT NULL,
    updated_at timestamp default current_timestamp NOT NULL,
    deleted_at timestamp,

    CONSTRAINT pk_portfolio PRIMARY KEY (user_id, account_id)
)