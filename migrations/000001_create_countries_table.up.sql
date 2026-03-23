CREATE TABLE IF NOT EXISTS countries (
    country_code CHAR(2)     NOT NULL,
    country_name VARCHAR(100) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_countries PRIMARY KEY (country_code)
);
