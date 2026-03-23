CREATE TABLE IF NOT EXISTS airports (
    airport_code CHAR(3)       NOT NULL,
    airport_name VARCHAR(150)  NOT NULL,
    city_id      INT           NOT NULL,
    latitude     NUMERIC(9, 6) NOT NULL,
    longitude    NUMERIC(9, 6) NOT NULL,
    timezone     VARCHAR(50)   NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_airports      PRIMARY KEY (airport_code),
    CONSTRAINT fk_airports_city FOREIGN KEY (city_id)
        REFERENCES cities (city_id) ON DELETE RESTRICT
);
