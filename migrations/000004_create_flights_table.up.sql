CREATE TABLE flights (
    id                      UUID         NOT NULL DEFAULT gen_random_uuid(),
    flight_number           VARCHAR(10)  NOT NULL,
    airline_id              UUID         NOT NULL,
    origin_airport_id       UUID         NOT NULL,
    destination_airport_id  UUID         NOT NULL,
    departure_time          TIMESTAMPTZ  NOT NULL,
    arrival_time            TIMESTAMPTZ  NOT NULL,
    duration_minutes        INTEGER      NOT NULL,
    stops                   SMALLINT     NOT NULL DEFAULT 0,
    aircraft_type           VARCHAR(50),
    total_seat_capacity     INTEGER      NOT NULL,
    status                  VARCHAR(20)  NOT NULL DEFAULT 'scheduled',
    estimated_departure_time TIMESTAMPTZ,
    estimated_arrival_time  TIMESTAMPTZ,
    departure_gate          VARCHAR(10),
    departure_terminal      VARCHAR(10),
    arrival_gate            VARCHAR(10),
    arrival_terminal        VARCHAR(10),
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT flights_pkey PRIMARY KEY (id),
    CONSTRAINT flights_airline_id_fkey
        FOREIGN KEY (airline_id) REFERENCES airlines (id),
    CONSTRAINT flights_origin_airport_id_fkey
        FOREIGN KEY (origin_airport_id) REFERENCES airports (id),
    CONSTRAINT flights_destination_airport_id_fkey
        FOREIGN KEY (destination_airport_id) REFERENCES airports (id),
    CONSTRAINT flights_duration_minutes_check
        CHECK (duration_minutes > 0),
    CONSTRAINT flights_stops_check
        CHECK (stops >= 0),
    CONSTRAINT flights_total_seat_capacity_check
        CHECK (total_seat_capacity > 0),
    CONSTRAINT flights_status_check
        CHECK (status IN ('scheduled','boarding','departed','en_route','landed','arrived','delayed','cancelled'))
);

CREATE INDEX flights_route_departure_idx
    ON flights (origin_airport_id, destination_airport_id, departure_time);
CREATE INDEX flights_airline_id_idx    ON flights (airline_id);
CREATE INDEX flights_flight_number_idx ON flights (flight_number);
CREATE INDEX flights_departure_time_idx ON flights (departure_time);
CREATE INDEX flights_status_idx        ON flights (status);
