CREATE TABLE IF NOT EXISTS availableitems (
    id bigserial NOT NULL PRIMARY KEY,
    knownitems_id bigint REFERENCES knownitems(id) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    expiration_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    container_size int NOT NULL,
    version integer NOT NULL DEFAULT 1
);

 /*
    expiration_at can be null possibly
 */