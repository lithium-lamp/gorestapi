package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"householdingindex.homecatalogue.net/internal/validator"
)

type KnownItemModel struct {
	DB *sql.DB
}

func (ki KnownItemModel) Insert(knownitem *KnownItem) error {
	query := `
		INSERT INTO knownitems (serial_number, long_name, short_name, tags, item_type, measurement, container_size)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, version`

	args := []interface{}{knownitem.SerialNumber, knownitem.LongName, knownitem.ShortName, pq.Array(knownitem.Tags), knownitem.ItemType, knownitem.Measurement, knownitem.ContainerSize}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return ki.DB.QueryRowContext(ctx, query, args...).Scan(&knownitem.ID, &knownitem.CreatedAt, &knownitem.Version)
}

func (ki KnownItemModel) Get(id int64) (*KnownItem, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, serial_number, long_name, short_name, tags, item_type, measurement, container_size, version
		FROM knownitems
		WHERE id = $1`

	var knownitem KnownItem

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := ki.DB.QueryRowContext(ctx, query, id).Scan(
		&knownitem.ID,
		&knownitem.CreatedAt,
		&knownitem.SerialNumber,
		&knownitem.LongName,
		&knownitem.ShortName,
		pq.Array(&knownitem.Tags),
		&knownitem.ItemType,
		&knownitem.Measurement,
		&knownitem.ContainerSize,
		&knownitem.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &knownitem, nil
}

func (ki KnownItemModel) GetAll(serialnumber int, longname string, shortname string, tags []string, itemtype int, measurement int, containersize int, filters Filters) ([]*KnownItem, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, serial_number, long_name, short_name, tags, item_type, measurement, container_size, version
		FROM knownitems
		WHERE (serial_number = $1 OR $1 = 0)
		AND (to_tsvector('simple', long_name) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (STRPOS(LOWER(short_name), LOWER($3)) > 0 OR $3 = '')
		AND (tags @> $4 OR $4 = '{}')
		AND (item_type = $5 OR $5 = 0)
		AND (measurement = $6 OR $6 = 0)
		AND (container_size = $7 OR $7 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $8 OFFSET $9`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{serialnumber, longname, shortname, pq.Array(tags), itemtype, measurement, containersize, filters.limit(), filters.offset()}

	rows, err := ki.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	knownitems := []*KnownItem{}

	for rows.Next() {
		var knownitem KnownItem

		err := rows.Scan(
			&totalRecords,
			&knownitem.ID,
			&knownitem.CreatedAt,
			&knownitem.SerialNumber,
			&knownitem.LongName,
			&knownitem.ShortName,
			pq.Array(&knownitem.Tags),
			&knownitem.ItemType,
			&knownitem.Measurement,
			&knownitem.ContainerSize,
			&knownitem.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		knownitems = append(knownitems, &knownitem)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return knownitems, metadata, nil
}

func (ki KnownItemModel) Update(knownitem *KnownItem) error {
	query := `
		UPDATE knownitems
		SET serial_number = $1, long_name = $2, short_name = $3, tags = $4, item_type = $5, measurement = $6, container_size = $7, version = version + 1
		WHERE id = $8 AND version = $9
		RETURNING version`

	args := []interface{}{
		knownitem.SerialNumber,
		knownitem.LongName,
		knownitem.ShortName,
		pq.Array(knownitem.Tags),
		knownitem.ItemType,
		knownitem.Measurement,
		knownitem.ContainerSize,
		knownitem.ID,
		knownitem.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := ki.DB.QueryRowContext(ctx, query, args...).Scan(&knownitem.Version)
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

func (ki KnownItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM knownitems
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := ki.DB.ExecContext(ctx, query, id)
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

type KnownItem struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	SerialNumber  int64     `json:"serial_number"`
	LongName      string    `json:"long_name"`
	ShortName     string    `json:"short_name"`
	Tags          []string  `json:"tags"`
	ItemType      int64     `json:"item_type"`
	Measurement   int64     `json:"measurement"`
	ContainerSize int32     `json:"container_size"`
	Version       int32     `json:"version"`
}

func ValidateKnownItem(v *validator.Validator, knownitem *KnownItem) {
	v.Check(knownitem.LongName != "", "long_name", "must be provided")
	v.Check(len(knownitem.LongName) <= 500, "long_name", "must not be more than 500 bytes long")

	v.Check(knownitem.ShortName != "", "short_name", "must be provided")
	v.Check(len(knownitem.ShortName) <= 100, "short_name", "must not be more than 100 bytes long")

	//v.Check(knownitem.Tags != nil, "tags", "must be provided")
	//v.Check(len(knownitem.Tags) >= 1, "tags", "must contain at least 1 tags")
	v.Check(len(knownitem.Tags) <= 100, "tags", "must not contain more than 100 tags")

	v.Check(validator.Unique(knownitem.Tags), "tags", "must not contain duplicate values")

	v.Check(knownitem.ItemType != 0, "item_type", "must be provided")
	v.Check(knownitem.ItemType >= 1, "item_type", "must be greater than 0")
	v.Check(knownitem.ItemType <= 6, "item_type", "must not be greater than 6") //TEMP VALUE

	v.Check(knownitem.Measurement != 0, "measurement", "must be provided")
	v.Check(knownitem.Measurement >= 1, "measurement", "must be greater than 0")
	v.Check(knownitem.Measurement <= 6, "measurement", "must not be greater than 6")

	v.Check(knownitem.ContainerSize >= 0, "container_size", "must be at least 0")
	v.Check(knownitem.ContainerSize <= 100000, "container_size", "must not be more than 100000 units")
}
