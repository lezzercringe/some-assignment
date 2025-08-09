-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id                TEXT PRIMARY KEY,
    track_number      TEXT NOT NULL,
    entry             TEXT NOT NULL,
    locale            TEXT NOT NULL,
    internal_signature TEXT NOT NULL,
    customer_id       TEXT NOT NULL,
    delivery_service  TEXT NOT NULL,
    shardkey          TEXT NOT NULL,
    sm_id             INTEGER NOT NULL,
    date_created      TIMESTAMPTZ NOT NULL,
    oof_shard         TEXT NOT NULL,
    delivery_name TEXT NOT NULL,
    delivery_city TEXT NOT NULL,
    delivery_phone TEXT NOT NULL,
    delivery_zip TEXT NOT NULL,
    delivery_address TEXT NOT NULL,
    delivery_region TEXT NOT NULL,
    delivery_email TEXT NOT NULL
); 

CREATE TABLE IF NOT EXISTS items (
    id          SERIAL PRIMARY KEY,
    order_id  TEXT REFERENCES orders(id) NOT NULL,
    chrt_id     INTEGER NOT NULL,
    track_number TEXT NOT NULL,
    price       DECIMAL NOT NULL,
    rid         TEXT NOT NULL,
    name        TEXT NOT NULL,
    sale        DECIMAL NOT NULL,
    size        TEXT NOT NULL,
    total_price DECIMAL NOT NULL,
    nm_id       INTEGER NOT NULL,
    brand       TEXT NOT NULL,
    status      INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
    transaction   TEXT PRIMARY KEY,
    order_id      TEXT UNIQUE NOT NULL REFERENCES orders(id),
    request_id    TEXT NOT NULL,
    currency      TEXT NOT NULL,
    provider      TEXT NOT NULL,
    amount        DECIMAL NOT NULL,
    payment_dt    BIGINT NOT NULL,
    bank          TEXT NOT NULL,
    delivery_cost DECIMAL NOT NULL,
    goods_total   INTEGER NOT NULL,
    custom_fee    DECIMAL NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_orders_date_created_desc ON orders (date_created DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_date_created_desc;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
