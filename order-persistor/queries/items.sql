-- name: CreateItem :one
INSERT INTO items(
    order_id,
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetItemsByOrderID :many
SELECT *
FROM items
WHERE order_id = $1;
