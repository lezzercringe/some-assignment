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

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;

-- name: GetRecentOrders :many
SELECT *
FROM orders
ORDER BY date_created DESC
LIMIT $1;
