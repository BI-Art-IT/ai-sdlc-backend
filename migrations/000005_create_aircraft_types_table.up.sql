CREATE TABLE IF NOT EXISTS aircraft_types (
    aircraft_type_id   SERIAL       NOT NULL,
    manufacturer       VARCHAR(100) NOT NULL,
    model              VARCHAR(100) NOT NULL,
    total_seats        INT          NOT NULL,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_aircraft_types              PRIMARY KEY (aircraft_type_id),
    CONSTRAINT uq_aircraft_types              UNIQUE (manufacturer, model),
    CONSTRAINT chk_aircraft_types_total_seats CHECK (total_seats > 0)
);
