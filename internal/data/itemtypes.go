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

type ItemTypeModel struct {
	DB *sql.DB
}

func (it ItemTypeModel) Insert(itemtype *ItemType) error {
	query := `
		INSERT INTO itemtypes (name)
		VALUES ($1)
		RETURNING id, created_at, version`

	args := []interface{}{itemtype.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return it.DB.QueryRowContext(ctx, query, args...).Scan(&itemtype.ID, &itemtype.CreatedAt, &itemtype.Version)
}

func (it ItemTypeModel) Get(id int64) (*ItemType, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, version
		FROM itemtypes
		WHERE id = $1`

	var itemtype ItemType

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := it.DB.QueryRowContext(ctx, query, id).Scan(
		&itemtype.ID,
		&itemtype.CreatedAt,
		&itemtype.Name,
		&itemtype.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &itemtype, nil
}

func (it ItemTypeModel) GetAll(name string, filters Filters) ([]*ItemType, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, name, version
		FROM itemtypes
		WHERE (name = $1 OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

	rows, err := it.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	itemtypes := []*ItemType{}

	for rows.Next() {
		var itemtype ItemType

		err := rows.Scan(
			&totalRecords,
			&itemtype.ID,
			&itemtype.CreatedAt,
			&itemtype.Name,
			&itemtype.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		itemtypes = append(itemtypes, &itemtype)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return itemtypes, metadata, nil
}

func (it ItemTypeModel) Update(itemtype *ItemType) error {
	query := `
		UPDATE itemtypes
		SET name = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version`

	args := []interface{}{
		itemtype.Name,
		itemtype.ID,
		itemtype.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := it.DB.QueryRowContext(ctx, query, args...).Scan(&itemtype.Version)
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

func (it ItemTypeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM itemtypes
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := it.DB.ExecContext(ctx, query, id)
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

type ItemType struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Version   int32     `json:"version"`
}

func ValidateItemType(v *validator.Validator, itemtype *ItemType) {
	v.Check(itemtype.Name != "", "name", "must be provided")
	v.Check(len(itemtype.Name) <= 500, "name", "must not be more than 500 bytes long")
	//v.Check(validator.Unique(input.Name), "name", "must not contain duplicate values")
}
