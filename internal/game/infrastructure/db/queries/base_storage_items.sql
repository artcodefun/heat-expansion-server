-- Base storage items queries

-- name: ListBaseStorageItems :many
SELECT id, base_id, prototype_id, status,
       present_data, state, created_at
FROM game.base_storage_items
WHERE base_id = @base_id
ORDER BY id;

-- name: DeleteBaseStorageItemsByBase :exec
DELETE FROM game.base_storage_items WHERE base_id = @base_id;

-- name: InsertBaseStorageItem :one
INSERT INTO game.base_storage_items (
    id, base_id, prototype_id, status,
    present_data, state, created_at
) VALUES (
    @id, @base_id, @prototype_id, @status,
    @present_data, @state, @created_at
)
RETURNING id;
