-- Army prototypes queries

-- name: GetArmyPrototypeByID :one
SELECT *
FROM game.army_item_prototypes
WHERE id = @id;

-- name: ListArmyPrototypes :many
SELECT *
FROM game.army_item_prototypes
ORDER BY id;

-- name: CreateArmyPrototype :exec
INSERT INTO game.army_item_prototypes (
    id, name, category, faction, unlock_technology_id,
    short_description, full_description, price, production_time,
    space, image_url, attack, defence, capacity, stealth, speed,
    creation_sources
) VALUES (
    @id, @name, @category, @faction, @unlock_technology_id,
    @short_description, @full_description, @price, @production_time,
    @space, @image_url, @attack, @defence, @capacity, @stealth, @speed,
    @creation_sources
);

-- name: UpdateArmyPrototype :one
UPDATE game.army_item_prototypes
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
    attack               = @attack,
    defence              = @defence,
    capacity             = @capacity,
    stealth              = @stealth,
    speed                = @speed,
    creation_sources     = @creation_sources
WHERE id = @id
RETURNING *;
