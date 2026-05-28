-- Storage item prototypes queries

-- name: GetStoragePrototypeByID :one
SELECT *
FROM game.storage_item_prototypes
WHERE id = @id;

-- name: ListStoragePrototypes :many
SELECT *
FROM game.storage_item_prototypes
ORDER BY id;
