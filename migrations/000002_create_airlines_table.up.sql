CREATE TABLE IF NOT EXISTS airlines (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    iata_code CHAR(2) NOT NULL,
    name VARCHAR(255) NOT NULL,
    logo_url VARCHAR(500),
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT airlines_iata_code_unique UNIQUE (iata_code)
);

CREATE INDEX IF NOT EXISTS airlines_iata_code_idx ON airlines (iata_code);
