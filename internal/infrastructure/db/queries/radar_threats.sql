-- name: InsertRadarThreat :one
INSERT INTO game.radar_threats (
    id, operation_id, owner_base_id, detected_at, detected_x, detected_y,
    source_x, source_y, target_x, target_y,
    estimated_arrival_at, arrival_at, type, status, attack, speed, stealth, capacity
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
) RETURNING *;

-- name: UpdateRadarThreat :one
UPDATE game.radar_threats
SET estimated_arrival_at = $2,
    arrival_at = $3,
    status = $4
WHERE id = $1
RETURNING *;

-- name: GetRadarThreat :one
SELECT * FROM game.radar_threats WHERE id = $1 LIMIT 1;

-- name: GetRadarThreatByOperationID :one
SELECT * FROM game.radar_threats WHERE operation_id = $1 LIMIT 1;

-- name: RadarThreatExists :one
SELECT EXISTS (
    SELECT 1 FROM game.radar_threats WHERE owner_base_id = $1 AND operation_id = $2
);
