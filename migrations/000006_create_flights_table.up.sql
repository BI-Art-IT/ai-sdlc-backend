-- departure_time / arrival_time represent the scheduled wall-clock time for
-- a recurring flight route.  The actual UTC datetime is derived at query time
-- by combining the travel date with these times and the respective airport
-- timezone (stored in the airports table).
CREATE TABLE IF NOT EXISTS flights (
    flight_id             BIGSERIAL    NOT NULL,
    airline_code          CHAR(2)      NOT NULL,
    flight_number         VARCHAR(10)  NOT NULL,
    origin_airport_code   CHAR(3)      NOT NULL,
    destination_airport_code CHAR(3)   NOT NULL,
    departure_time        TIME         NOT NULL,
    arrival_time          TIME         NOT NULL,
    days_of_week          SMALLINT[]   NOT NULL,
    aircraft_type_id      INT          NOT NULL,
    base_price            NUMERIC(12, 2) NOT NULL,
    available_seats       INT          NOT NULL,
    status                VARCHAR(20)  NOT NULL DEFAULT 'active',
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_flights                    PRIMARY KEY (flight_id),
    CONSTRAINT uq_flights_number             UNIQUE (airline_code, flight_number),
    CONSTRAINT fk_flights_airline            FOREIGN KEY (airline_code)
        REFERENCES airlines (airline_code) ON DELETE RESTRICT,
    CONSTRAINT fk_flights_origin             FOREIGN KEY (origin_airport_code)
        REFERENCES airports (airport_code) ON DELETE RESTRICT,
    CONSTRAINT fk_flights_destination        FOREIGN KEY (destination_airport_code)
        REFERENCES airports (airport_code) ON DELETE RESTRICT,
    CONSTRAINT fk_flights_aircraft_type      FOREIGN KEY (aircraft_type_id)
        REFERENCES aircraft_types (aircraft_type_id) ON DELETE RESTRICT,
    CONSTRAINT chk_flights_status            CHECK (status IN ('active', 'cancelled', 'suspended')),
    CONSTRAINT chk_flights_available_seats   CHECK (available_seats >= 0),
    CONSTRAINT chk_flights_base_price        CHECK (base_price >= 0),
    CONSTRAINT chk_flights_different_airports CHECK (origin_airport_code <> destination_airport_code)
);

CREATE INDEX IF NOT EXISTS idx_flights_origin_destination
    ON flights (origin_airport_code, destination_airport_code);

CREATE INDEX IF NOT EXISTS idx_flights_airline
    ON flights (airline_code);

CREATE INDEX IF NOT EXISTS idx_flights_departure_time
    ON flights (departure_time);
