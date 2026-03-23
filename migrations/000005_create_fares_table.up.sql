CREATE TABLE fares (
    id              UUID           NOT NULL DEFAULT gen_random_uuid(),
    flight_id       UUID           NOT NULL,
    fare_class_id   UUID           NOT NULL,
    base_fare       DECIMAL(10,2)  NOT NULL,
    taxes           DECIMAL(10,2)  NOT NULL,
    fees            DECIMAL(10,2)  NOT NULL,
    total_price     DECIMAL(10,2)  NOT NULL,
    currency        CHAR(3)        NOT NULL,
    available_seats SMALLINT       NOT NULL,
    valid_from      TIMESTAMPTZ    NOT NULL,
    valid_until     TIMESTAMPTZ    NOT NULL,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fares_pkey PRIMARY KEY (id),
    CONSTRAINT fares_flight_id_fkey
        FOREIGN KEY (flight_id) REFERENCES flights (id),
    CONSTRAINT fares_fare_class_id_fkey
        FOREIGN KEY (fare_class_id) REFERENCES fare_classes (id)
);

CREATE INDEX fares_flight_fare_class_idx ON fares (flight_id, fare_class_id);
CREATE INDEX fares_valid_until_idx       ON fares (valid_until);
CREATE INDEX fares_total_price_idx       ON fares (total_price);
