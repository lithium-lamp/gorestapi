package data

import (
	"time"

	"householdingindex.homecatalogue.net/internal/validator"
)

type AvailableItem struct {
	ID            int64       `json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	ExpirationAt  time.Time   `json:"expiration_at,omitempty"`
	LongName      string      `json:"long_name"`
	ShortName     string      `json:"short_name"`
	ItemType      ItemType    `json:"item_type"`
	Measurement   Measurement `json:"measurement"`
	ContainerSize int32       `json:"container_size"`
}

func ValidateAvailableItem(v *validator.Validator, availableitem *AvailableItem) {
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
}
