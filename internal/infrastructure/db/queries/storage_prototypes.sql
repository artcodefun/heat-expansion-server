-- Storage item prototypes queries

-- name: GetStoragePrototypeByID :one
SELECT id, name, category, short_description, full_description, image_url,
       buff_data, map_data, damaged_data, artifact_data, consumable_data
FROM storage_item_prototypes
WHERE id = @id;

-- name: ListStoragePrototypes :many
SELECT id, name, category, short_description, full_description, image_url,
       buff_data, map_data, damaged_data, artifact_data, consumable_data
FROM storage_item_prototypes
ORDER BY id;
