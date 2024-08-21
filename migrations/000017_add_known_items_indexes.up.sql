CREATE INDEX IF NOT EXISTS knownitems_long_name_idx ON knownitems USING GIN (to_tsvector('simple', long_name));

CREATE INDEX IF NOT EXISTS knownitems_short_name_idx ON knownitems USING GIN (to_tsvector('simple', short_name));

CREATE INDEX IF NOT EXISTS knownitems_tags_idx ON knownitems USING GIN (tags);