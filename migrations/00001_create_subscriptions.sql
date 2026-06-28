-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_month DATE NOT NULL,
    end_month DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_user_id
    ON subscriptions(user_id);

CREATE INDEX idx_subscriptions_service_name
    ON subscriptions(service_name);

CREATE INDEX idx_subscriptions_period
    ON subscriptions(start_month, end_month);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd