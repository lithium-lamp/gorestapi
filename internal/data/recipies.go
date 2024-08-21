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

type RecipeModel struct {
	DB *sql.DB
}

func (rm RecipeModel) Insert(recipe *Recipe) error {
	query := `
		INSERT INTO recipies (name, description, cooking_steps, cook_time_minutes, portions, tags)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version`

	args := []interface{}{recipe.Name, recipe.Description, pq.Array(recipe.CookingSteps), recipe.CookTimeMinutes, recipe.Portions, pq.Array(recipe.Tags)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rm.DB.QueryRowContext(ctx, query, args...).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.Version)
}

func (rm RecipeModel) Get(id int64) (*Recipe, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, description, cooking_steps, cook_time_minutes, portions, tags, version
		FROM recipies
		WHERE id = $1`

	var recipe Recipe

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := rm.DB.QueryRowContext(ctx, query, id).Scan(
		&recipe.ID,
		&recipe.CreatedAt,
		&recipe.Name,
		&recipe.Description,
		pq.Array(&recipe.CookingSteps),
		&recipe.CookTimeMinutes,
		&recipe.Portions,
		pq.Array(&recipe.Tags),
		&recipe.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipe, nil
}

func (rm RecipeModel) GetAll(name string, cooktimeminutes int, portions int, tags []string, filters Filters) ([]*Recipe, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, name, description, cooking_steps, cook_time_minutes, portions, tags, version
		FROM recipies
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (cook_time_minutes = $2 OR $2 = 0)
		AND (portions = $3 OR $3 = 0)
		AND (tags @> $4 OR $4 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, cooktimeminutes, portions, pq.Array(tags), filters.limit(), filters.offset()}

	rows, err := rm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	recipies := []*Recipe{}

	for rows.Next() {
		var recipe Recipe

		err := rows.Scan(
			&totalRecords,
			&recipe.ID,
			&recipe.CreatedAt,
			&recipe.Name,
			&recipe.CookTimeMinutes,
			&recipe.Portions,
			pq.Array(&recipe.Tags),
			&recipe.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		recipies = append(recipies, &recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return recipies, metadata, nil
}

func (rm RecipeModel) Update(recipe *Recipe) error {
	query := `
		UPDATE recipies
		SET name = $1, description = $2, cooking_steps = $3, cook_time_minutes = $4, portions = $5, tags = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version`

	args := []interface{}{
		recipe.Name,
		recipe.Description,
		pq.Array(recipe.CookingSteps),
		recipe.CookTimeMinutes,
		recipe.Portions,
		pq.Array(recipe.Tags),
		recipe.ID,
		recipe.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rm.DB.QueryRowContext(ctx, query, args...).Scan(&recipe.Version)
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

func (rm RecipeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM recipies
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := rm.DB.ExecContext(ctx, query, id)
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

type Recipe struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CookingSteps    []string  `json:"cooking_steps"`
	CookTimeMinutes int32     `json:"cook_time_minutes"`
	Portions        int32     `json:"portions"`
	Tags            []string  `json:"tags"`
	Version         int32     `json:"version"`
}

func ValidateRecipe(v *validator.Validator, recipe *Recipe) {
	v.Check(recipe.Name != "", "name", "must be provided")
	v.Check(len(recipe.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(recipe.Description != "", "description", "must be provided")
	v.Check(len(recipe.Description) <= 5000, "description", "must not be more than 5000 bytes long")

	v.Check(len(recipe.Tags) <= 100, "cooking_steps", "must not contain more than 100 cooking steps")

	v.Check(recipe.CookTimeMinutes != 0, "cook_time_minutes", "must be provided")
	v.Check(recipe.CookTimeMinutes >= 1, "cook_time_minutes", "must be greater than 0")
	v.Check(recipe.CookTimeMinutes <= 10000, "cook_time_minutes", "must not be greater than 10000") //TEMP VALUE

	v.Check(recipe.Portions != 0, "portions", "must be provided")
	v.Check(recipe.Portions >= 1, "portions", "must be greater than 0")
	v.Check(recipe.Portions <= 10000, "portions", "must not be greater than 10000") //TEMP VALUE

	//v.Check(knownitem.Tags != nil, "tags", "must be provided")
	//v.Check(len(knownitem.Tags) >= 1, "tags", "must contain at least 1 tags")
	v.Check(len(recipe.Tags) <= 100, "tags", "must not contain more than 100 tags")

	v.Check(validator.Unique(recipe.Tags), "tags", "must not contain duplicate values")
}
