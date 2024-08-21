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

type IngredientModel struct {
	DB *sql.DB
}

func (im IngredientModel) Insert(ingredient *Ingredient) error {
	query := `
		INSERT INTO ingredients (name, tags)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []interface{}{ingredient.Name, pq.Array(ingredient.Tags)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return im.DB.QueryRowContext(ctx, query, args...).Scan(&ingredient.ID, &ingredient.CreatedAt, &ingredient.Version)
}

func (im IngredientModel) Get(id int64) (*Ingredient, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, tags, version
		FROM ingredients
		WHERE id = $1`

	var ingredient Ingredient

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := im.DB.QueryRowContext(ctx, query, id).Scan(
		&ingredient.ID,
		&ingredient.CreatedAt,
		&ingredient.Name,
		pq.Array(&ingredient.Tags),
		&ingredient.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &ingredient, nil
}

func (im IngredientModel) GetAll(name string, tags []string, filters Filters) ([]*Ingredient, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, name, tags, version
		FROM ingredients
		WHERE (name = $1 OR $1 = '')
		AND (tags @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, pq.Array(tags), filters.limit(), filters.offset()}

	rows, err := im.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	ingredients := []*Ingredient{}

	for rows.Next() {
		var ingredient Ingredient

		err := rows.Scan(
			&totalRecords,
			&ingredient.ID,
			&ingredient.CreatedAt,
			&ingredient.Name,
			pq.Array(&ingredient.Tags),
			&ingredient.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		ingredients = append(ingredients, &ingredient)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return ingredients, metadata, nil
}

func (im IngredientModel) Update(ingredient *Ingredient) error {
	query := `
		UPDATE ingredients
		SET name = $1, tags = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []interface{}{
		ingredient.Name,
		pq.Array(ingredient.Tags),
		ingredient.ID,
		ingredient.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := im.DB.QueryRowContext(ctx, query, args...).Scan(&ingredient.Version)
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

func (im IngredientModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM ingredients
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := im.DB.ExecContext(ctx, query, id)
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

type Ingredient struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Tags      []string  `json:"tags"`
	Version   int32     `json:"version"`
}

func ValidateIngredient(v *validator.Validator, ingredient *Ingredient) {
	v.Check(ingredient.Name != "", "name", "must be provided")
	v.Check(len(ingredient.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(len(ingredient.Tags) <= 100, "tags", "must not contain more than 100 tags")
	v.Check(validator.Unique(ingredient.Tags), "tags", "must not contain duplicate values")
}
