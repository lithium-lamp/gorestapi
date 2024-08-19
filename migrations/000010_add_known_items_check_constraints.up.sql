ALTER TABLE knownitems ADD CONSTRAINT knownitems_itemtype_check CHECK (item_type >= 1);

ALTER TABLE knownitems ADD CONSTRAINT knownitems_measurement_check CHECK (measurement >= 1);

ALTER TABLE knownitems ADD CONSTRAINT knownitems_tags_check CHECK (array_length(genres, 1) BETWEEN 1 AND 100);
/*

    buissiness logic can be added here
*/