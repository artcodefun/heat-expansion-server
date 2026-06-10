-- Technology prototypes queries

-- name: GetTechPrototypeByID :one
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, improvement
FROM game.tech_item_prototypes
WHERE id = @id;

-- name: ListTechPrototypes :many
SELECT id, name, category, unlock_technology_id, short_description, full_description,
       price,
       research_time, image_url, improvement
FROM game.tech_item_prototypes
ORDER BY id;

-- name: CreateTechPrototype :exec
INSERT INTO game.tech_item_prototypes (
    id, name, category, unlock_technology_id, short_description, full_description,
    price, research_time, image_url, improvement
) VALUES (
    @id, @name, @category, @unlock_technology_id, @short_description, @full_description,
    @price, @research_time, @image_url, @improvement
);

-- name: UpdateTechPrototype :one
UPDATE game.tech_item_prototypes SET
    name                 = @name,
    category             = @category,
    unlock_technology_id = @unlock_technology_id,
    short_description    = @short_description,
    full_description     = @full_description,
    price                = @price,
    research_time        = @research_time,
    image_url            = @image_url,
    improvement          = @improvement
WHERE id = @id
RETURNING *;
