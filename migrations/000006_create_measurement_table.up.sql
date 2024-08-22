CREATE TABLE IF NOT EXISTS measurements (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL UNIQUE,
    short_name text NOT NULL UNIQUE,
    version integer NOT NULL DEFAULT 1
);

INSERT INTO measurements (name, short_name)
VALUES
    ('units', 'x'),
    ('grams', 'g'),
    ('milliliters', 'ml');