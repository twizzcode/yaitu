ALTER TABLE sessions
ADD COLUMN hostname TEXT,
ADD COLUMN family_id UUID NOT NULL DEFAULT gen_random_uuid();

ALTER TABLE sessions
ADD CONSTRAINT sessions_hostname_valid CHECK (
    hostname IS NULL
    OR (
        hostname = LOWER(hostname)
        AND hostname NOT LIKE '%/%'
        AND hostname NOT LIKE '%:%'
        AND hostname NOT LIKE '.%'
        AND hostname NOT LIKE '%.'
        AND hostname <> ''
    )
);

CREATE INDEX sessions_family_id_idx
ON sessions(family_id);

CREATE INDEX sessions_hostname_idx
ON sessions(hostname);