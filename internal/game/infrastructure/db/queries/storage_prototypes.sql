-- Storage item prototypes queries

-- name: GetStoragePrototypeByID :one
SELECT *
FROM game.storage_item_prototypes
WHERE id = @id;

-- name: ListStoragePrototypes :many
SELECT *
FROM game.storage_item_prototypes
ORDER BY id;

-- name: CreateStoragePrototype :exec
INSERT INTO game.storage_item_prototypes (
    id, name, category, estimated_worth,
    short_description, full_description, image_url,
    buff_data, intel_data, damaged_data, artifact_data, consumable_data,
    creation_sources
) VALUES (
    @id, @name, @category, @estimated_worth,
    @short_description, @full_description, @image_url,
    @buff_data, @intel_data, @damaged_data, @artifact_data, @consumable_data,
    @creation_sources
);

-- name: UpdateStoragePrototype :one
UPDATE game.storage_item_prototypes SET
    name               = @name,
    category           = @category,
    estimated_worth    = @estimated_worth,
    short_description  = @short_description,
    full_description   = @full_description,
    image_url          = @image_url,
    buff_data          = @buff_data,
    intel_data         = @intel_data,
    damaged_data       = @damaged_data,
    artifact_data      = @artifact_data,
    consumable_data    = @consumable_data,
    creation_sources   = @creation_sources
WHERE id = @id
RETURNING *;
