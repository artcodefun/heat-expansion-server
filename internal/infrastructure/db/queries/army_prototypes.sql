-- Army prototypes queries

-- name: GetArmyPrototypeByID :one
SELECT id, name, category, faction, unlock_technology_id, short_description, full_description,
       price,
       production_time, space, image_url,
       attack, defence, capacity, stealth, speed
FROM game.army_item_prototypes
WHERE id = @id;

-- name: ListArmyPrototypes :many
SELECT id, name, category, faction, unlock_technology_id, short_description, full_description,
       price,
       production_time, space, image_url,
       attack, defence, capacity, stealth, speed
FROM game.army_item_prototypes
ORDER BY id;
