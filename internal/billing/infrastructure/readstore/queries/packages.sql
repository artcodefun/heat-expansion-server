-- name: ListActivePackages :many
SELECT id, crystals, price_minor_units, currency, image_url
FROM billing.crystal_packages
WHERE is_active = TRUE
ORDER BY price_minor_units ASC;

-- name: ListAllPackages :many
SELECT id, name, crystals, price_minor_units, currency, image_url, is_active
FROM billing.crystal_packages
ORDER BY price_minor_units ASC;

-- name: GetPackageByID :one
SELECT id, name, crystals, price_minor_units, currency, image_url, is_active
FROM billing.crystal_packages
WHERE id = @id;
