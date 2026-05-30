-- name: ListActivePackages :many
SELECT id, name, crystals, price_minor_units, currency, image_url
FROM billing.crystal_packages
WHERE is_active = TRUE
ORDER BY price_minor_units ASC;
