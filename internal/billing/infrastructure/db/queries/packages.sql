-- name: GetPackageByID :one
SELECT id, name, crystals, price_minor_units, currency, image_url, is_active, created_at, updated_at
FROM billing.crystal_packages
WHERE id = $1;

-- name: CreatePackage :one
INSERT INTO billing.crystal_packages (id, name, crystals, price_minor_units, currency, image_url, is_active, created_at, updated_at)
VALUES (@id, @name, @crystals, @price_minor_units, @currency, @image_url, @is_active, @created_at, @updated_at)
RETURNING *;

-- name: UpdatePackage :one
UPDATE billing.crystal_packages
SET name              = @name,
    crystals          = @crystals,
    price_minor_units = @price_minor_units,
    currency          = @currency,
    image_url         = @image_url,
    is_active         = @is_active,
    updated_at        = @updated_at
WHERE id = @id
RETURNING *;
