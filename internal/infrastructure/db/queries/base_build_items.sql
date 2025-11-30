-- Base build items queries

-- name: ListBaseBuildItems :many
SELECT id, base_id, prototype_id, status,
    pending_data, in_prod_data, present_data,
       created_at
FROM base_build_items
WHERE base_id = @base_id
ORDER BY id;

-- name: DeleteBaseBuildItemsByBase :exec
DELETE FROM base_build_items WHERE base_id = @base_id;

-- name: InsertBaseBuildItem :one
INSERT INTO base_build_items (
    base_id, prototype_id, status,
    pending_data, in_prod_data, present_data,
    created_at
) VALUES (
    @base_id, @prototype_id, @status,
    @pending_data, @in_prod_data, @present_data,
    @created_at
)
RETURNING id;
