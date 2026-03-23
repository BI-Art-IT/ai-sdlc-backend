CREATE TABLE airlines (
    id        UUID         NOT NULL DEFAULT gen_random_uuid(),
    iata_code CHAR(2)      NOT NULL,
    name      VARCHAR(255) NOT NULL,
    logo_url  VARCHAR(500),
    active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT airlines_pkey PRIMARY KEY (id),
    CONSTRAINT airlines_iata_code_key UNIQUE (iata_code)
);
