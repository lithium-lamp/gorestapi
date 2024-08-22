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

type TagModel struct {
	DB *sql.DB
}

func (tm TagModel) Insert(tag *Tag) error {
	query := `
		INSERT INTO tags (itemtype, name)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []interface{}{tag.ItemType, tag.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return tm.DB.QueryRowContext(ctx, query, args...).Scan(&tag.ID, &tag.CreatedAt, &tag.Version)
}

func (tm TagModel) Get(id int64) (*Tag, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, itemtype, name, version
		FROM tags
		WHERE id = $1`

	var tag Tag

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := tm.DB.QueryRowContext(ctx, query, id).Scan(
		&tag.ID,
		&tag.CreatedAt,
		&tag.ItemType,
		&tag.Name,
		&tag.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tag, nil
}

func (tm TagModel) GetAll(itemtype int, name string, filters Filters) ([]*Tag, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, itemtype, name, version
		FROM tags
		WHERE (itemtype = $1 OR $1 = 0)
		AND (name = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{itemtype, name, filters.limit(), filters.offset()}

	rows, err := tm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	tags := []*Tag{}

	for rows.Next() {
		var tag Tag

		err := rows.Scan(
			&totalRecords,
			&tag.ID,
			&tag.CreatedAt,
			&tag.ItemType,
			&tag.Name,
			&tag.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tags, metadata, nil
}

func (tm TagModel) Update(tag *Tag) error {
	query := `
		UPDATE measurements
		SET itemtype = $1, name = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []interface{}{
		tag.ItemType,
		tag.Name,
		tag.ID,
		tag.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tm.DB.QueryRowContext(ctx, query, args...).Scan(&tag.Version)
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

func (tm TagModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM tags
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := tm.DB.ExecContext(ctx, query, id)
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

type Tag struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ItemType  int64     `json:"itemtype"`
	Name      string    `json:"name"`
	Version   int32     `json:"version"`
}

func ValidateTag(v *validator.Validator, tag *Tag) {
	v.Check(tag.ItemType != 0, "itemtype", "must be provided")
	v.Check(tag.ItemType >= 1, "itemtype", "must be greater than 0")
	v.Check(tag.ItemType <= 6, "itemtype", "must not be greater than 6") //TEMP VALUE

	v.Check(tag.Name != "", "name", "must be provided")
	v.Check(len(tag.Name) <= 500, "name", "must not be more than 500 bytes long")
}

/*
	id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    itemtype bigint REFERENCES itemtypes(id) NOT NULL,
    name text NOT NULL UNIQUE,
    version integer NOT NULL DEFAULT 1


*/
