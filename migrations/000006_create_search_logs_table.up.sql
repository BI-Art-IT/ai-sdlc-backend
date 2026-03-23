CREATE TABLE IF NOT EXISTS search_logs (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id_hash CHAR(64),
    session_id VARCHAR(64),
    origin_code VARCHAR(3) NOT NULL,
    destination_code VARCHAR(3) NOT NULL,
    departure_date DATE NOT NULL,
    return_date DATE,
    trip_type VARCHAR(20) NOT NULL,
    adults SMALLINT NOT NULL,
    children SMALLINT NOT NULL,
    infants SMALLINT NOT NULL,
    cabin_class VARCHAR(20) NOT NULL,
    flexible_dates BOOLEAN NOT NULL DEFAULT false,
    results_returned INTEGER,
    cache_hit BOOLEAN,
    response_time_ms INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS search_logs_created_at_idx ON search_logs (created_at);
CREATE INDEX IF NOT EXISTS search_logs_origin_dest_idx ON search_logs (origin_code, destination_code);
