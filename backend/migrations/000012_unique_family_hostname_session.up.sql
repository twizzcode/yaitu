WITH ranked_sessions AS (
    SELECT
        ctid,
        ROW_NUMBER() OVER (
            PARTITION BY family_id, hostname
            ORDER BY created_at DESC
        ) AS row_number
    FROM sessions
    WHERE hostname IS NOT NULL
)
DELETE FROM sessions
USING ranked_sessions
WHERE sessions.ctid = ranked_sessions.ctid
  AND ranked_sessions.row_number > 1;

CREATE UNIQUE INDEX sessions_family_hostname_idx
ON sessions(family_id, hostname)
WHERE hostname IS NOT NULL;