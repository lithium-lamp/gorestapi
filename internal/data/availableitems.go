package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"householdingindex.homecatalogue.net/internal/validator"
)

type AvailableItemModel struct {
	DB *sql.DB
}

func (ai AvailableItemModel) Insert(availableitem *AvailableItem) error {
	query := `
		INSERT INTO availableitems (expiration_at, long_name, short_name, item_type, measurement, container_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version`

	args := []interface{}{availableitem.ExpirationAt, availableitem.LongName, availableitem.ShortName, availableitem.ItemType, availableitem.Measurement, availableitem.ContainerSize}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return ai.DB.QueryRowContext(ctx, query, args...).Scan(&availableitem.ID, &availableitem.CreatedAt, &availableitem.Version)
}

func (ai AvailableItemModel) Get(id int64) (*AvailableItem, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, expiration_at, long_name, short_name, item_type, measurement, container_size, version
		FROM availableitems
		WHERE id = $1`

	var availableitem AvailableItem

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := ai.DB.QueryRowContext(ctx, query, id).Scan(
		&availableitem.ID,
		&availableitem.CreatedAt,
		&availableitem.ExpirationAt,
		&availableitem.LongName,
		&availableitem.ShortName,
		&availableitem.ItemType,
		&availableitem.Measurement,
		&availableitem.ContainerSize,
		&availableitem.Version,
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

func (ai AvailableItemModel) GetAll(expirationat time.Time, longname string, shortname string, itemtype int, measurement Measurement, containersize int, filters Filters) ([]*AvailableItem, Metadata, error) {
	//expiration_at currently retrieves items larger than the input ====> search for items that are still fresh according to current date
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, expiration_at, long_name, short_name, item_type, measurement, container_size, version
		FROM availableitems
		WHERE (expiration_at >= $1 OR $1 = '0001-01-01T00:00:00Z')
		AND (to_tsvector('simple', long_name) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (STRPOS(LOWER(short_name), LOWER($3)) > 0 OR $3 = '')
		AND (item_type = $4 OR $4 = 0)
		AND (measurement = $5 OR $5 = 0)
		AND (container_size = $6 OR $6 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $7 OFFSET $8`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{expirationat.Format(time.RFC3339), longname, shortname, itemtype, measurement, containersize, filters.limit(), filters.offset()}

	rows, err := ai.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	availableitems := []*AvailableItem{}

	for rows.Next() {
		var availableitem AvailableItem

		err := rows.Scan(
			&totalRecords,
			&availableitem.ID,
			&availableitem.CreatedAt,
			&availableitem.ExpirationAt,
			&availableitem.LongName,
			&availableitem.ShortName,
			&availableitem.ItemType,
			&availableitem.Measurement,
			&availableitem.ContainerSize,
			&availableitem.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		availableitems = append(availableitems, &availableitem)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return availableitems, metadata, nil
}

func (ai AvailableItemModel) Update(availableitem *AvailableItem) error {
	query := `
		UPDATE availableitems
		SET expiration_at = $1, long_name = $2, short_name = $3, item_type = $4, measurement = $5, container_size = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version`

	args := []interface{}{
		availableitem.ExpirationAt,
		availableitem.LongName,
		availableitem.ShortName,
		availableitem.ItemType,
		availableitem.Measurement,
		availableitem.ContainerSize,
		availableitem.ID,
		availableitem.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := ai.DB.QueryRowContext(ctx, query, args...).Scan(&availableitem.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (ai AvailableItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM availableitems
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := ai.DB.ExecContext(ctx, query, id)
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
	ItemType      int64       `json:"item_type"`
	Measurement   Measurement `json:"measurement"`
	ContainerSize int32       `json:"container_size"`
	Version       int32       `json:"version"`
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
