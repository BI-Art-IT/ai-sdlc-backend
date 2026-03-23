CREATE TABLE IF NOT EXISTS cities (
    city_id      SERIAL       NOT NULL,
    city_name    VARCHAR(100) NOT NULL,
    country_code CHAR(2)      NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_cities         PRIMARY KEY (city_id),
    CONSTRAINT fk_cities_country FOREIGN KEY (country_code)
        REFERENCES countries (country_code) ON DELETE RESTRICT
);
