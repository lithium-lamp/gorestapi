ALTER TABLE availableitems ADD CONSTRAINT availableitems_itemtype_check CHECK (item_type >= 1);

ALTER TABLE availableitems ADD CONSTRAINT availableitems_measurement_check CHECK (measurement >= 1);

/*
    ALTER TABLE availableitems ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
*/

/*

    v.Check(availableitem.LongName != "", "long_name", "must be provided")
	v.Check(len(availableitem.LongName) <= 500, "long_name", "must not be more than 500 bytes long")
	//v.Check(validator.Unique(input.LongName), "long_name", "must not contain duplicate values")

	v.Check(availableitem.ShortName != "", "short_name", "must be provided")
	v.Check(len(availableitem.ShortName) <= 100, "short_name", "must not be more than 100 bytes long")

	v.Check(availableitem.ItemType != 0, "item_type", "must be provided")
	v.Check(availableitem.ItemType >= 1, "item_type", "must be greater than 0")
	v.Check(availableitem.ItemType <= 6, "item_type", "must not be greater than 6") //TEMP VALUE

	v.Check(availableitem.Measurement != 0, "measurement", "must be provided")
	v.Check(availableitem.Measurement >= 1, "measurement", "must be greater than 0")
	v.Check(availableitem.Measurement <= 6, "measurement", "must not be greater than 6")

	v.Check(availableitem.ContainerSize >= 0, "container_size", "must be at least 0")
	v.Check(availableitem.ContainerSize <= 100000, "container_size", "must not be more than 100000 units")

*/