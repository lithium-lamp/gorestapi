CREATE TABLE IF NOT EXISTS knownitems (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    serial_number bigint NOT NULL,
    long_name text NOT NULL,
    short_name text NOT NULL,
    tags text [] NOT NULL,
    item_type bigint REFERENCES itemtypes(id) NOT NULL,
    measurement bigint REFERENCES measurements(id) NOT NULL,
    container_size int NOT NULL,
    version integer NOT NULL DEFAULT 1
);