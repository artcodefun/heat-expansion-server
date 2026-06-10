-- Prototypes queries for read-store

-- name: ListArmyPrototypes :many
SELECT *
FROM game.army_item_prototypes
ORDER BY id;

-- name: GetArmyPrototypeByID :one
SELECT *
FROM game.army_item_prototypes
WHERE id = $1;

-- name: ListBuildPrototypes :many
SELECT *
FROM game.build_item_prototypes
ORDER BY id;

-- name: GetBuildPrototypeByID :one
SELECT *
FROM game.build_item_prototypes
WHERE id = $1;

-- name: GetStoragePrototypeByID :one
SELECT *
FROM game.storage_item_prototypes
WHERE id = $1;

-- name: ListStoragePrototypes :many
SELECT *
FROM game.storage_item_prototypes
ORDER BY id;

-- name: ListTechPrototypes :many
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, improvement
FROM game.tech_item_prototypes
ORDER BY id;

-- name: GetTechPrototypeByID :one
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, improvement
FROM game.tech_item_prototypes
WHERE id = $1;
