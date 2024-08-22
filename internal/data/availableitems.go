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
		INSERT INTO availableitems (knownitems_id, expiration_at, container_size)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []interface{}{availableitem.KnownItemsID, availableitem.ExpirationAt, availableitem.ContainerSize}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return ai.DB.QueryRowContext(ctx, query, args...).Scan(&availableitem.ID, &availableitem.CreatedAt, &availableitem.Version)
}

func (ai AvailableItemModel) Get(id int64) (*AvailableItem, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, knownitems_id, created_at, expiration_at, container_size, version
		FROM availableitems
		WHERE id = $1`

	var availableitem AvailableItem

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := ai.DB.QueryRowContext(ctx, query, id).Scan(
		&availableitem.ID,
		&availableitem.KnownItemsID,
		&availableitem.CreatedAt,
		&availableitem.ExpirationAt,
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

func (ai AvailableItemModel) GetAll(knownitemsid int, expirationat time.Time, containersize int, filters Filters) ([]*AvailableItem, Metadata, error) {
	//expiration_at currently retrieves items larger than the input ====> search for items that are still fresh according to current date
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, knownitems_id, created_at, expiration_at, container_size, version
		FROM availableitems
		WHERE (knownitems_id = $1 OR $1 = 0)
		AND (expiration_at >= $2 OR $2 = '0001-01-01T00:00:00Z')
		AND (container_size = $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{knownitemsid, expirationat.Format(time.RFC3339), containersize, filters.limit(), filters.offset()}

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
			&availableitem.KnownItemsID,
			&availableitem.CreatedAt,
			&availableitem.ExpirationAt,
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
		SET knownitems_id = $1, expiration_at = $2, container_size = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []interface{}{
		availableitem.KnownItemsID,
		availableitem.ExpirationAt,
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
	ID            int64     `json:"id"`
	KnownItemsID  int64     `json:"knownitems_id"`
	CreatedAt     time.Time `json:"created_at"`
	ExpirationAt  time.Time `json:"expiration_at,omitempty"`
	ContainerSize int32     `json:"container_size"`
	Version       int32     `json:"version"`
}

func ValidateAvailableItem(v *validator.Validator, availableitem *AvailableItem) {
	v.Check(availableitem.KnownItemsID >= 0, "knownitems_id", "must be at least 0")
	v.Check(availableitem.KnownItemsID <= 100000, "knownitems_id", "must not be more than 100000 units")

	v.Check(availableitem.ContainerSize >= 0, "container_size", "must be at least 0")
	v.Check(availableitem.ContainerSize <= 100000, "container_size", "must not be more than 100000 units")
}
