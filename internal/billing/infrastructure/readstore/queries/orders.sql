-- name: GetOrderByID :one
SELECT id, user_id, package_id, crystals, amount_minor_units, currency, provider, status, confirmation_url, created_at
FROM billing.purchase_orders
WHERE id = $1;
