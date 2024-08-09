ALTER TABLE availableitems ADD CONSTRAINT availableitems_itemtype_check CHECK (item_type >= 1);

ALTER TABLE availableitems ADD CONSTRAINT availableitems_measurement_check CHECK (measurement >= 1);

/*
    ALTER TABLE availableitems ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

    buissiness logic can be added here
*/