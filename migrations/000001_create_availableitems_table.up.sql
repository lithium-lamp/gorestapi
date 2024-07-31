CREATE TABLE IF NOT EXISTS availableitems (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    expiration_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    long_name text NOT NULL,
    short_name text NOT NULL,
    item_type int NOT NULL,
    measurement int NOT NULL,
    container_size int NOT NULL
);