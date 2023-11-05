-- name: GetCalendarByID :one
SELECT * FROM calendar WHERE name_key = ?;

-- name: UpsertCalendar :one
INSERT INTO calendar (
        name_key,
        created_at,
        updated_at
) VALUES (
        ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) ON CONFLICT (name_key) DO UPDATE SET
        updated_at = CURRENT_TIMESTAMP
RETURNING *;