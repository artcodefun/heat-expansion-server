-- Base tech items queries

-- name: ListBaseTechItems :many
SELECT id, base_id, prototype_id, status,
       in_progress_data, done_data,
       created_at
FROM base_tech_items
WHERE base_id = @base_id
ORDER BY id;

-- name: DeleteBaseTechItemsByBase :exec
DELETE FROM base_tech_items WHERE base_id = @base_id;

-- name: InsertBaseTechItem :one
INSERT INTO base_tech_items (
    id, base_id, prototype_id, status,
    in_progress_data, done_data,
    created_at
) VALUES (
    @id, @base_id, @prototype_id, @status,
    @in_progress_data, @done_data,
    @created_at
)
RETURNING id;
