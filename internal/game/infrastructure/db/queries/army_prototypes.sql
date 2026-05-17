-- Army prototypes queries

-- name: GetArmyPrototypeByID :one
SELECT *
FROM game.army_item_prototypes
WHERE id = @id;

-- name: ListArmyPrototypes :many
SELECT *
FROM game.army_item_prototypes
ORDER BY id;
