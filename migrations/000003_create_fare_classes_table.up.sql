CREATE TABLE fare_classes (
    id         UUID        NOT NULL DEFAULT gen_random_uuid(),
    code       VARCHAR(20) NOT NULL,
    name       VARCHAR(50) NOT NULL,
    sort_order SMALLINT    NOT NULL,

    CONSTRAINT fare_classes_pkey PRIMARY KEY (id),
    CONSTRAINT fare_classes_code_key UNIQUE (code)
);

INSERT INTO fare_classes (code, name, sort_order) VALUES
    ('economy',         'Economy',         1),
    ('premium-economy', 'Premium Economy', 2),
    ('business',        'Business',        3),
    ('first',           'First Class',     4);
