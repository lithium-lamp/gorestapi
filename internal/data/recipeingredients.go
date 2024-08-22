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

type RecipeIngredientModel struct {
	DB *sql.DB
}

func (rm RecipeIngredientModel) Insert(recipeingredient *RecipeIngredient) error {
	query := `
		INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, measurement)
		VALUES ($1, $2, $3, $4)
		RETURNING recipe_id, ingredient_id, created_at, version`

	args := []interface{}{recipeingredient.RecipeID, recipeingredient.IngredientID, recipeingredient.Amount, recipeingredient.Measurement}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return rm.DB.QueryRowContext(ctx, query, args...).Scan(&recipeingredient.RecipeID, &recipeingredient.IngredientID, &recipeingredient.CreatedAt, &recipeingredient.Version)
}

func (rm RecipeIngredientModel) Get(recipeid int64, ingredientid int64) (*RecipeIngredient, error) {
	if recipeid < 1 {
		return nil, ErrRecordNotFound
	}

	if ingredientid < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT recipe_id, ingredient_id, created_at, amount, measurement, version
		FROM recipe_ingredients
		WHERE recipe_id = $1 AND ingredient_id = $2`

	var recipeingredient RecipeIngredient

	args := []interface{}{recipeid, ingredientid}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rm.DB.QueryRowContext(ctx, query, args...).Scan(
		&recipeingredient.RecipeID,
		&recipeingredient.IngredientID,
		&recipeingredient.CreatedAt,
		&recipeingredient.Amount,
		&recipeingredient.Measurement,
		&recipeingredient.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipeingredient, nil
}

func (rm RecipeIngredientModel) GetAll(amount int, measurement int, filters Filters) ([]*RecipeIngredient, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), recipe_id, ingredient_id, created_at, amount, measurement, version
		FROM recipe_ingredients
		WHERE (amount = $1 OR $1 = 0)
		AND (measurement = $2 OR $2 = 0)
		ORDER BY %s %s, recipe_id ASC, ingredient_id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{amount, measurement, filters.limit(), filters.offset()}

	rows, err := rm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	recipeingredients := []*RecipeIngredient{}

	for rows.Next() {
		var recipeingredient RecipeIngredient

		err := rows.Scan(
			&totalRecords,
			&recipeingredient.RecipeID,
			&recipeingredient.IngredientID,
			&recipeingredient.CreatedAt,
			&recipeingredient.Amount,
			&recipeingredient.Measurement,
			&recipeingredient.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		recipeingredients = append(recipeingredients, &recipeingredient)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return recipeingredients, metadata, nil
}

func (mm RecipeIngredientModel) Update(recipeingredient *RecipeIngredient, oldrecipeid *int64, oldingredientid *int64) error {
	query := `
		UPDATE recipe_ingredients
		SET recipe_id = $1, ingredient_id = $2, amount = $3, measurement = $4, version = version + 1
		WHERE recipe_id = $5 AND ingredient_id = $6 AND version = $7
		RETURNING version`

	args := []interface{}{
		recipeingredient.RecipeID,
		recipeingredient.IngredientID,
		recipeingredient.Amount,
		recipeingredient.Measurement,
		oldrecipeid,
		oldingredientid,
		recipeingredient.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mm.DB.QueryRowContext(ctx, query, args...).Scan(&recipeingredient.Version)
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

func (rm RecipeIngredientModel) Delete(recipeid int64, ingredientid int64) error {
	if recipeid < 1 {
		return ErrRecordNotFound
	}

	if ingredientid < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM recipe_ingredients
		WHERE recipe_id = $1 AND ingredient_id = $2`

	args := []interface{}{recipeid, ingredientid}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := rm.DB.ExecContext(ctx, query, args...)
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

type RecipeIngredient struct {
	RecipeID     int64     `json:"recipe_id"`
	IngredientID int64     `json:"ingredient_id"`
	CreatedAt    time.Time `json:"created_at"`
	Amount       int32     `json:"amount"`
	Measurement  int64     `json:"measurement"`
	Version      int32     `json:"version"`
}

func ValidateRecipeIngredient(v *validator.Validator, recipeingredient *RecipeIngredient) {
	v.Check(recipeingredient.RecipeID != 0, "recipe_id", "must be provided")
	v.Check(recipeingredient.RecipeID >= 1, "recipe_id", "must be at least 1")

	v.Check(recipeingredient.IngredientID != 0, "ingredient_id", "must be provided")
	v.Check(recipeingredient.IngredientID >= 1, "ingredient_id", "must be at least 1")

	v.Check(recipeingredient.Amount != 0, "amount", "must be provided")
	v.Check(recipeingredient.Amount >= 1, "amount", "must be at least 1")
	v.Check(recipeingredient.Amount <= 100000, "amount", "must not be larger than 100000")

	v.Check(recipeingredient.Measurement != 0, "measurement", "must be provided")
	v.Check(recipeingredient.Measurement >= 1, "measurement", "must be at least 1")
}
