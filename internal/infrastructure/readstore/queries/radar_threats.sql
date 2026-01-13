-- name: GetRadarThreat :one
SELECT * FROM radar_threats WHERE id = $1 LIMIT 1;

-- name: ListIncomingThreats :many
SELECT * FROM radar_threats WHERE owner_base_id = $1 AND status = 'ARRIVING' ORDER BY estimated_arrival_at ASC;
