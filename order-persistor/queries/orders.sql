-- name: CreateOrder :one
INSERT INTO orders (
    id,
    track_number,
    entry,
    locale,
    internal_signature,
    customer_id,
    delivery_service,
    shardkey,
    sm_id,
    date_created,
    oof_shard,
    delivery_name,
    delivery_phone,
    delivery_zip,
    delivery_address,
    delivery_region,
    delivery_email,
    delivery_city
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET
    track_number = $2,
    entry = $3,
    locale = $4,
    internal_signature = $5,
    customer_id = $6,
    delivery_service = $7,
    shardkey = $8,
    sm_id = $9,
    date_created = $10,
    oof_shard = $11,
    delivery_name = $12,
    delivery_phone = $13,
    delivery_zip = $14,
    delivery_address = $15,
    delivery_region = $16,
    delivery_email = $17,
    delivery_city = $18
WHERE id = $1
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;
