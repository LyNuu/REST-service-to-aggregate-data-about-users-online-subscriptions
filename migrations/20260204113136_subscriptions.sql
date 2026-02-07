-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions
(
    id           UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id      UUID         NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price        INT          NOT NULL CHECK (price >= 0),
    start_date   DATE         NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX idx_subscriptions_period ON subscriptions (start_date, end_date);

INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
SELECT (ARRAY [
    '60601fee-2bf1-4721-ae6f-7636e79a0cba',
    '50af7296-8c0d-4eef-b894-a08b4367b5ef',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    '84b86377-d9f3-424a-93f8-580a17441551'
    ])[floor(random() * 4 + 1)]::UUID                                                                                                                   as user_id,

       (ARRAY ['Netflix', 'Spotify', 'Yandex Plus', 'YouTube Premium', 'PS Plus', 'Xbox Live', 'Apple Music', 'ChatGPT Plus'])[floor(random() * 8 + 1)] as service_name,

       (floor(random() * 19 + 1) * 100)::INT                                                                                                            as price,

       start_dt                                                                                                                                         as start_date,

       (start_dt + (floor(random() * 12 + 1) * INTERVAL '1 month'))::DATE                                                                               as end_date
FROM (SELECT CURRENT_DATE - (random() * 730)::INT as start_dt
      FROM generate_series(1, 100)) as subquery;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
