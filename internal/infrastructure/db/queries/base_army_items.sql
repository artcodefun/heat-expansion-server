-- Base army items queries

-- name: ListBaseArmyItems :many
SELECT id, base_id, prototype_id, status,
       pending_data, in_prod_data, present_data, deployed_data,
       created_at
FROM base_army_items
WHERE base_id = @base_id
ORDER BY id;

-- name: DeleteBaseArmyItemsByBase :exec
DELETE FROM base_army_items WHERE base_id = @base_id;

-- name: InsertBaseArmyItem :one
INSERT INTO base_army_items (
    id, base_id, prototype_id, status,
    pending_data, in_prod_data, present_data, deployed_data,
    created_at
) VALUES (
    @id, @base_id, @prototype_id, @status,
    @pending_data, @in_prod_data, @present_data, @deployed_data,
    @created_at
)
RETURNING id;
