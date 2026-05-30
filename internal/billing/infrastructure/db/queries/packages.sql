-- name: GetPackageByID :one
SELECT id, name, crystals, price_minor_units, currency, image_url, is_active, created_at, updated_at
FROM billing.crystal_packages
WHERE id = $1;
