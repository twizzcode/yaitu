CREATE TABLE venue_domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id UUID NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    hostname TEXT NOT NULL,
    type TEXT NOT NULL
        CHECK (type IN ('platform', 'custom')),
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'active', 'failed')),
    verification_token_hash BYTEA,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (hostname = LOWER(hostname)),
    CHECK (hostname NOT LIKE '%/%'),
    CHECK (hostname NOT LIKE '%:%'),
    CHECK (hostname NOT LIKE '.%'),
    CHECK (hostname NOT LIKE '%.'),
    CHECK (hostname <> ''),
    CHECK (
        status <> 'active'
        OR verified_at IS NOT NULL
    )
);

CREATE UNIQUE INDEX venue_domains_hostname_idx
ON venue_domains(hostname);

CREATE UNIQUE INDEX venue_domains_platform_venue_idx
ON venue_domains(venue_id)
WHERE type = 'platform';

CREATE INDEX venue_domains_venue_id_idx
ON venue_domains(venue_id);

CREATE INDEX venue_domains_status_idx
ON venue_domains(status);