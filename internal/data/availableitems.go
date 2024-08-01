package data

import (
	"database/sql"
	"errors"
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
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, expiration_at, long_name, short_name, item_type, measurement, container_size
		FROM availableitems
		WHERE id = $1`

	var availableitem AvailableItem

	err := ai.DB.QueryRow(query, id).Scan(
		&availableitem.ID,
		&availableitem.CreatedAt,
		&availableitem.ExpirationAt,
		&availableitem.LongName,
		&availableitem.ShortName,
		&availableitem.ItemType,
		&availableitem.Measurement,
		&availableitem.ContainerSize,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &availableitem, nil
}

func (ai AvailableItemModel) Update(availableitem *AvailableItem) error {
	query := `
		UPDATE availableitems
		SET expiration_at = $1, long_name = $2, short_name = $3, item_type = $4, measurement = $5, container_size = $6
		WHERE id = $7
		RETURNING id`

	args := []interface{}{
		availableitem.ExpirationAt,
		availableitem.LongName,
		availableitem.ShortName,
		availableitem.ItemType,
		availableitem.Measurement,
		availableitem.ContainerSize,
		availableitem.ID,
	}

	return ai.DB.QueryRow(query, args...).Scan(&availableitem.ID)
}

func (ai AvailableItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM availableitems WHERE id = $1`

	result, err := ai.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

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
