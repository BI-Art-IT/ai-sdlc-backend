CREATE TABLE IF NOT EXISTS airports (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    iata_code CHAR(3) NOT NULL,
    name VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    country_code CHAR(2) NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    latitude DECIMAL(9, 6),
    longitude DECIMAL(9, 6),
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT airports_iata_code_unique UNIQUE (iata_code)
);

CREATE INDEX IF NOT EXISTS airports_iata_code_idx ON airports (iata_code);
CREATE INDEX IF NOT EXISTS airports_city_idx ON airports (city);
CREATE INDEX IF NOT EXISTS airports_country_code_idx ON airports (country_code);
