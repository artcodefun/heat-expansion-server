-- Building prototypes queries

-- name: GetBuildPrototypeByID :one
SELECT *
FROM game.build_item_prototypes
WHERE id = @id;

-- name: ListBuildPrototypes :many
SELECT *
FROM game.build_item_prototypes
ORDER BY id;

-- name: CreateBuildPrototype :exec
INSERT INTO game.build_item_prototypes (
    id, name, category, faction, unlock_technology_id,
    short_description, full_description, price, production_time,
    space, image_url, control_data, resources_data, defense_data,
    military_data, intelligence_data, creation_sources
) VALUES (
    @id, @name, @category, @faction, @unlock_technology_id,
    @short_description, @full_description, @price, @production_time,
    @space, @image_url, @control_data, @resources_data, @defense_data,
    @military_data, @intelligence_data, @creation_sources
);

-- name: UpdateBuildPrototype :one
UPDATE game.build_item_prototypes
SET name                 = @name,
    category             = @category,
    faction              = @faction,
    unlock_technology_id = @unlock_technology_id,
    short_description    = @short_description,
    full_description     = @full_description,
    price                = @price,
    production_time      = @production_time,
    space                = @space,
    image_url            = @image_url,
    control_data         = @control_data,
    resources_data       = @resources_data,
    defense_data         = @defense_data,
    military_data        = @military_data,
    intelligence_data    = @intelligence_data,
    creation_sources     = @creation_sources
WHERE id = @id
RETURNING *;
