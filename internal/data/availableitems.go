package data

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"householdingindex.homecatalogue.net/internal/validator"
)

type AvailableItemModel struct {
	DB *sql.DB
}

func (ai AvailableItemModel) Insert(availableitem *AvailableItem) error {
	query := `
		INSERT INTO availableitems (long_name, short_name, item_type, measurement, container_size)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, expiration_at`

	args := []interface{}{availableitem.LongName, availableitem.ShortName, availableitem.ItemType, availableitem.Measurement, availableitem.ContainerSize}

	return ai.DB.QueryRow(query, args...).Scan(&availableitem.ID, &availableitem.CreatedAt, &availableitem.ExpirationAt)
}

func (ai AvailableItemModel) Get(id int64) (*AvailableItem, error) {
	return nil, nil
}

func (ai AvailableItemModel) Update(availableitem *AvailableItem) error {
	return nil
}

func (ai AvailableItemModel) Delete(id int64) error {
	return nil
}

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
