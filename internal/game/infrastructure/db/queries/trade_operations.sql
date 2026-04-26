-- Trade operations queries

-- name: GetTradeOperationByID :one
SELECT *
FROM game.trade_operations
WHERE id = @id;

-- name: GetTradeOperationByIDForUpdate :one
SELECT *
FROM game.trade_operations
WHERE id = @id
FOR UPDATE;

-- name: InsertTradeOperation :one
INSERT INTO game.trade_operations (
    operation_uuid,
    sender_user_id,
    sender_base_id,
    receiver_user_id,
    receiver_base_id,
    created_at,
    source_x,
    source_y,
    target_x,
    target_y,
    offered_payload,
    requested_payload,
    transport_units,
    storage_snaps,
    total_modifiers,
    expires_at,
    outbound_depart_at,
    outbound_arrive_at,
    arrived_at_target_at,
    return_depart_at,
    return_arrive_at,
    completed_at,
    phase,
    result,
    crystals_skip_price
) VALUES (
    @operation_uuid,
    @sender_user_id,
    @sender_base_id,
    @receiver_user_id,
    @receiver_base_id,
    @created_at,
    @source_x,
    @source_y,
    @target_x,
    @target_y,
    @offered_payload,
    @requested_payload,
    @transport_units,
    @storage_snaps,
    @total_modifiers,
    @expires_at,
    @outbound_depart_at,
    @outbound_arrive_at,
    @arrived_at_target_at,
    @return_depart_at,
    @return_arrive_at,
    @completed_at,
    @phase,
    @result,
    @crystals_skip_price
)
RETURNING id;

-- name: UpdateTradeOperation :exec
UPDATE game.trade_operations
SET operation_uuid = @operation_uuid,
    created_at = @created_at,
    sender_user_id = @sender_user_id,
    sender_base_id = @sender_base_id,
    receiver_user_id = @receiver_user_id,
    receiver_base_id = @receiver_base_id,
    source_x = @source_x,
    source_y = @source_y,
    target_x = @target_x,
    target_y = @target_y,
    offered_payload = @offered_payload,
    requested_payload = @requested_payload,
    transport_units = @transport_units,
    storage_snaps = @storage_snaps,
    total_modifiers = @total_modifiers,
    expires_at = @expires_at,
    outbound_depart_at = @outbound_depart_at,
    outbound_arrive_at = @outbound_arrive_at,
    arrived_at_target_at = @arrived_at_target_at,
    return_depart_at = @return_depart_at,
    return_arrive_at = @return_arrive_at,
    completed_at = @completed_at,
    phase = @phase,
    result = @result,
    crystals_skip_price = @crystals_skip_price
WHERE id = @id;

-- name: DeleteTradeOperation :exec
DELETE FROM game.trade_operations
WHERE id = @id;