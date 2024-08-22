ALTER TABLE knownitems DROP CONSTRAINT IF EXISTS knownitems_itemtype_check;

ALTER TABLE knownitems DROP CONSTRAINT IF EXISTS knownitems_measurement_check;

ALTER TABLE knownitems DROP CONSTRAINT IF EXISTS knownitems_tags_check;

/*
    ALTER TABLE knownitems DROP CONSTRAINT IF EXISTS genres_length_check;
*/