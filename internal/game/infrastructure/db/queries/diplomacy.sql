-- name: InsertDiplomaticRelationship :exec
INSERT INTO game.diplomatic_relationships (
    id,
    user_a_id,
    user_b_id,
    status,
    changed_by_user_id,
    changed_at,
    war_declared_at,
    war_attacks_allowed_at,
    neutrality_protected_until
  ) VALUES (@id, @user_a_id, @user_b_id, @status, @changed_by_user_id, @changed_at, @war_declared_at, @war_attacks_allowed_at, @neutrality_protected_until);

-- name: UpdateDiplomaticRelationship :exec
UPDATE game.diplomatic_relationships
  SET status = @status,
    changed_by_user_id = @changed_by_user_id,
    changed_at = @changed_at,
    war_declared_at = @war_declared_at,
    war_attacks_allowed_at = @war_attacks_allowed_at,
    neutrality_protected_until = @neutrality_protected_until
  WHERE user_a_id = @user_a_id AND user_b_id = @user_b_id;

-- name: GetDiplomaticRelationship :one
SELECT * FROM game.diplomatic_relationships
  WHERE user_a_id = @user_a_id AND user_b_id = @user_b_id;

-- name: GetDiplomaticRelationshipForUpdate :one
SELECT * FROM game.diplomatic_relationships
  WHERE user_a_id = @user_a_id AND user_b_id = @user_b_id
FOR UPDATE;

-- name: InsertDiplomaticMessage :exec
INSERT INTO game.diplomatic_messages (
    id,
    sender_user_id,
    receiver_user_id,
    sender_base_id,
    receiver_base_id,
    request_id,
    reply_to_message_id,
    is_read,
    content,
    created_at
) VALUES (@id, @sender_user_id, @receiver_user_id, @sender_base_id, @receiver_base_id, @request_id, @reply_to_message_id, @is_read, @content, @created_at);

-- name: ExistsDiplomaticMessageByRequestAndContent :one
SELECT EXISTS (
  SELECT 1
  FROM game.diplomatic_messages
  WHERE request_id = @request_id
    AND content = @content
);

-- name: GetDiplomaticMessageByRequestAndContent :one
SELECT * FROM game.diplomatic_messages
WHERE request_id = @request_id
  AND content = @content
ORDER BY created_at ASC
LIMIT 1;

-- name: GetDiplomaticMessage :one
SELECT * FROM game.diplomatic_messages
WHERE id = @id;

-- name: MarkDiplomaticChatAsRead :exec
UPDATE game.diplomatic_messages
SET is_read = TRUE
WHERE receiver_user_id = @receiver_user_id
  AND sender_user_id = @sender_user_id
  AND is_read = FALSE;

-- name: InsertDiplomaticRequest :exec
INSERT INTO game.diplomatic_requests (
    id,
    sender_user_id,
    receiver_user_id,
    sender_base_id,
    receiver_base_id,
    kind,
    status,
    created_at,
    resolved_at,
    expires_at
) VALUES (@id, @sender_user_id, @receiver_user_id, @sender_base_id, @receiver_base_id, @kind, @status, @created_at, @resolved_at, @expires_at);

-- name: UpdateDiplomaticRequest :exec
UPDATE game.diplomatic_requests
SET status = @status,
    resolved_at = @resolved_at
WHERE id = @id;

-- name: GetDiplomaticRequest :one
SELECT * FROM game.diplomatic_requests
WHERE id = @id;

-- name: GetDiplomaticRequestForUpdate :one
SELECT * FROM game.diplomatic_requests
WHERE id = @id
FOR UPDATE;

-- name: ExistsPendingDiplomaticRequestByKind :one
SELECT EXISTS (
    SELECT 1
    FROM game.diplomatic_requests
    WHERE sender_user_id IN (@user_a_id, @user_b_id)
      AND receiver_user_id IN (@user_a_id, @user_b_id)
      AND status = 'PENDING'
      AND kind = @kind
);
