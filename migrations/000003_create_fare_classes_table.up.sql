CREATE TABLE IF NOT EXISTS fare_classes (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    sort_order SMALLINT NOT NULL,
    CONSTRAINT fare_classes_code_unique UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS fare_classes_code_idx ON fare_classes (code);

INSERT INTO fare_classes (code, name, sort_order) VALUES
    ('economy', 'Economy', 1),
    ('premium-economy', 'Premium Economy', 2),
    ('business', 'Business', 3),
    ('first', 'First Class', 4)
ON CONFLICT (code) DO NOTHING;
