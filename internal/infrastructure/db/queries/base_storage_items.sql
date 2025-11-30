-- Base storage items queries

-- name: ListBaseStorageItems :many
SELECT id, base_id, prototype_id, status,
       present_data, state, created_at
FROM base_storage_items
WHERE base_id = @base_id
ORDER BY id;

-- name: DeleteBaseStorageItemsByBase :exec
DELETE FROM base_storage_items WHERE base_id = @base_id;

-- name: InsertBaseStorageItem :one
INSERT INTO base_storage_items (
    base_id, prototype_id, status,
    present_data, state, created_at
) VALUES (
    @base_id, @prototype_id, @status,
    @present_data, @state, @created_at
)
RETURNING id;
