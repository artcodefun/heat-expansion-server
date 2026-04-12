-- name: ListDiplomaticRelationships :many
SELECT
    CASE WHEN dr.user_a_id = @current_user_id THEN dr.user_b_id ELSE dr.user_a_id END AS other_user_id,
    u.name AS other_user_name,
    dr.status,
    dr.changed_by_user_id,
    dr.changed_at,
    dr.war_declared_at,
    dr.war_attacks_allowed_at,
    dr.neutrality_protected_until
FROM game.diplomatic_relationships dr
JOIN game.users u ON u.id = CASE WHEN dr.user_a_id = @current_user_id THEN dr.user_b_id ELSE dr.user_a_id END
WHERE (dr.user_a_id = @current_user_id OR dr.user_b_id = @current_user_id)
    AND (@status = '' OR dr.status = @status)
ORDER BY dr.changed_at DESC;

-- name: GetDiplomaticRelationship :one
SELECT
    CASE WHEN dr.user_a_id = @current_user_id THEN dr.user_b_id ELSE dr.user_a_id END AS other_user_id,
    u.name AS other_user_name,
    dr.status,
    dr.changed_by_user_id,
    dr.changed_at,
    dr.war_declared_at,
    dr.war_attacks_allowed_at,
    dr.neutrality_protected_until
FROM game.diplomatic_relationships dr
JOIN game.users u ON u.id = CASE WHEN dr.user_a_id = @current_user_id THEN dr.user_b_id ELSE dr.user_a_id END
WHERE (dr.user_a_id = @current_user_id AND dr.user_b_id = @other_user_id)
   OR (dr.user_a_id = @other_user_id AND dr.user_b_id = @current_user_id);

-- name: GetDiplomaticRequest :one
SELECT
    dr.id,
    dr.sender_user_id,
    su.name AS sender_user_name,
    dr.receiver_user_id,
    ru.name AS receiver_user_name,
    dr.sender_base_id,
    dr.receiver_base_id,
    dr.kind,
    dr.status,
    dr.created_at,
    dr.resolved_at,
    dr.expires_at
FROM game.diplomatic_requests dr
JOIN game.users su ON su.id = dr.sender_user_id
JOIN game.users ru ON ru.id = dr.receiver_user_id
WHERE dr.id = @id;

-- name: ListDiplomaticChats :many
WITH user_messages AS (
    SELECT
        dm.id,
        dm.sender_user_id,
        dm.receiver_user_id,
        dm.sender_base_id,
        dm.receiver_base_id,
        dm.request_id,
        dm.reply_to_message_id,
        dm.is_read,
        dm.content,
        dm.created_at,
        CASE WHEN dm.sender_user_id = @current_user_id THEN dm.receiver_user_id ELSE dm.sender_user_id END AS other_user_id,
        ROW_NUMBER() OVER (
            PARTITION BY CASE WHEN dm.sender_user_id = @current_user_id THEN dm.receiver_user_id ELSE dm.sender_user_id END
            ORDER BY dm.created_at DESC
        ) AS row_num
    FROM game.diplomatic_messages dm
    WHERE dm.sender_user_id = @current_user_id OR dm.receiver_user_id = @current_user_id
), unread_counts AS (
    SELECT sender_user_id AS other_user_id, COUNT(*) AS unread_count
    FROM game.diplomatic_messages
    WHERE receiver_user_id = @current_user_id AND is_read = FALSE
    GROUP BY sender_user_id
)
SELECT
    um.other_user_id,
    u.name AS other_user_name,
    um.id,
    um.sender_user_id,
    su.name AS sender_user_name,
    um.receiver_user_id,
    ru.name AS receiver_user_name,
    um.sender_base_id,
    um.receiver_base_id,
    um.request_id,
    um.reply_to_message_id,
    um.is_read,
    um.content,
    um.created_at,
    COALESCE(uc.unread_count, 0)::BIGINT AS unread_count
FROM user_messages um
JOIN game.users u ON u.id = um.other_user_id
JOIN game.users su ON su.id = um.sender_user_id
JOIN game.users ru ON ru.id = um.receiver_user_id
LEFT JOIN unread_counts uc ON uc.other_user_id = um.other_user_id
WHERE um.row_num = 1
ORDER BY um.created_at DESC;

-- name: CountUnreadDiplomaticMessagesByUser :one
SELECT count(*)
FROM game.diplomatic_messages
WHERE receiver_user_id = @receiver_user_id AND is_read = FALSE;

-- name: ListDiplomaticMessagesByChat :many
SELECT
    dm.id,
    dm.sender_user_id,
    su.name AS sender_user_name,
    dm.receiver_user_id,
    ru.name AS receiver_user_name,
    dm.sender_base_id,
    dm.receiver_base_id,
    dm.request_id,
    dm.reply_to_message_id,
    dm.is_read,
    dm.content,
    dm.created_at
FROM game.diplomatic_messages dm
JOIN game.users su ON su.id = dm.sender_user_id
JOIN game.users ru ON ru.id = dm.receiver_user_id
WHERE (dm.sender_user_id = @current_user_id AND dm.receiver_user_id = @other_user_id)
   OR (dm.sender_user_id = @other_user_id AND dm.receiver_user_id = @current_user_id)
ORDER BY dm.created_at ASC;

-- name: ListPendingDiplomaticRequests :many
SELECT
    dr.id,
    dr.sender_user_id,
    su.name AS sender_user_name,
    dr.receiver_user_id,
    ru.name AS receiver_user_name,
    dr.sender_base_id,
    dr.receiver_base_id,
    dr.kind,
    dr.status,
    dr.created_at,
    dr.resolved_at,
    dr.expires_at
FROM game.diplomatic_requests dr
JOIN game.users su ON su.id = dr.sender_user_id
JOIN game.users ru ON ru.id = dr.receiver_user_id
WHERE dr.receiver_user_id = @receiver_user_id AND dr.status = 'PENDING'
ORDER BY dr.created_at DESC;
