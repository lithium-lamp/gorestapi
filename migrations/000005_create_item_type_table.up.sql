CREATE TABLE IF NOT EXISTS itemtypes (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL UNIQUE,
    version integer NOT NULL DEFAULT 1
);

INSERT INTO itemtypes (name)
VALUES
    ('food'),
    ('household'),
    ('furniture'),
    ('vehicles');