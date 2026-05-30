-- name: InsertOrder :exec
INSERT INTO billing.purchase_orders (id, user_id, package_id, crystals, amount_minor_units, currency, provider, status, provider_order_id, confirmation_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: UpdateOrder :exec
UPDATE billing.purchase_orders
SET status = $2, provider_order_id = $3, confirmation_url = $4, updated_at = $5
WHERE id = $1;

-- name: GetOrderByID :one
SELECT id, user_id, package_id, crystals, amount_minor_units, currency, provider, status, provider_order_id, confirmation_url, created_at, updated_at
FROM billing.purchase_orders
WHERE id = $1;

-- name: GetOrderByProviderOrderID :one
SELECT id, user_id, package_id, crystals, amount_minor_units, currency, provider, status, provider_order_id, confirmation_url, created_at, updated_at
FROM billing.purchase_orders
WHERE provider_order_id = $1;

-- name: GetOrderByProviderOrderIDForUpdate :one
SELECT id, user_id, package_id, crystals, amount_minor_units, currency, provider, status, provider_order_id, confirmation_url, created_at, updated_at
FROM billing.purchase_orders
WHERE provider_order_id = $1
FOR UPDATE;
