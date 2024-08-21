CREATE INDEX IF NOT EXISTS recipies_name_idx ON recipies USING GIN (to_tsvector('simple', name));

CREATE INDEX IF NOT EXISTS recipies_tags_idx ON recipies USING GIN (tags);

CREATE INDEX IF NOT EXISTS ingredients_name_idx ON ingredients USING GIN (to_tsvector('simple', name));

CREATE INDEX IF NOT EXISTS ingredients_tags_idx ON ingredients USING GIN (tags);