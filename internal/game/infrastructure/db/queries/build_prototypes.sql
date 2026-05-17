-- Building prototypes queries

-- name: GetBuildPrototypeByID :one
SELECT *
FROM game.build_item_prototypes
WHERE id = @id;

-- name: ListBuildPrototypes :many
SELECT *
FROM game.build_item_prototypes
ORDER BY id;
