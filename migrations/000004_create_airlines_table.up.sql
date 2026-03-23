CREATE TABLE IF NOT EXISTS airlines (
    airline_code CHAR(2)      NOT NULL,
    airline_name VARCHAR(100) NOT NULL,
    country_code CHAR(2)      NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_airlines         PRIMARY KEY (airline_code),
    CONSTRAINT fk_airlines_country FOREIGN KEY (country_code)
        REFERENCES countries (country_code) ON DELETE RESTRICT
);
