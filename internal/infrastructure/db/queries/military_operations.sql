-- Military operations queries

-- name: GetMilitaryOperationByID :one
SELECT id, type, owner_user_id, source_base_id,
       source_x, source_y, target_x, target_y,
        outbound_depart_at, outbound_arrive_at,
        return_depart_at, return_arrive_at,
        completed_at, phase, result,
        units, spy_result, attack_result,
        crystals_skip_price
FROM military_operations
WHERE id = @id;

-- name: GetMilitaryOperationByIDForUpdate :one
SELECT id, type, owner_user_id, source_base_id,
       source_x, source_y, target_x, target_y,
        outbound_depart_at, outbound_arrive_at,
        return_depart_at, return_arrive_at,
        completed_at, phase, result,
        units, spy_result, attack_result,
        crystals_skip_price
FROM military_operations
WHERE id = @id
FOR UPDATE;

-- name: ListOpsBySourceBase :many
SELECT id, type, owner_user_id, source_base_id,
       source_x, source_y, target_x, target_y,
        outbound_depart_at, outbound_arrive_at,
        return_depart_at, return_arrive_at,
        completed_at, phase, result,
        units, spy_result, attack_result,
        crystals_skip_price
FROM military_operations
WHERE source_base_id = @base_id
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: ListOpsByTargetCoordinates :many
SELECT id, type, owner_user_id, source_base_id,
       source_x, source_y, target_x, target_y,
        outbound_depart_at, outbound_arrive_at,
        return_depart_at, return_arrive_at,
        completed_at, phase, result,
        units, spy_result, attack_result,
        crystals_skip_price
FROM military_operations
WHERE target_x = @target_x AND target_y = @target_y
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: InsertMilitaryOperation :one
INSERT INTO military_operations (
        type, owner_user_id, source_base_id,
    source_x, source_y, target_x, target_y,
    outbound_depart_at, outbound_arrive_at,
    return_depart_at, return_arrive_at,
        completed_at, phase, result,
        units, spy_result, attack_result,
        crystals_skip_price
) VALUES (
        @type, @owner_user_id, @source_base_id,
    @source_x, @source_y, @target_x, @target_y,
    @outbound_depart_at, @outbound_arrive_at,
    @return_depart_at, @return_arrive_at,
        @completed_at, @phase, @result,
        @units, @spy_result, @attack_result,
        @crystals_skip_price
)
RETURNING id;

-- name: UpdateMilitaryOperation :exec
UPDATE military_operations
SET type = @type,
        owner_user_id = @owner_user_id,
        source_base_id = @source_base_id,
        source_x = @source_x,
        source_y = @source_y,
        target_x = @target_x,
        target_y = @target_y,
        outbound_depart_at = @outbound_depart_at,
        outbound_arrive_at = @outbound_arrive_at,
        return_depart_at = @return_depart_at,
        return_arrive_at = @return_arrive_at,
        completed_at = @completed_at,
        phase = @phase,
        result = @result,
        units = @units,
        spy_result = @spy_result,
        attack_result = @attack_result,
        crystals_skip_price = @crystals_skip_price
WHERE id = @id;

-- name: DeleteMilitaryOperation :exec
DELETE FROM military_operations WHERE id = @id;
