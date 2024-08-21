package main

import (
	"errors"
	"fmt"
	"net/http"

	"householdingindex.homecatalogue.net/internal/data"
	"householdingindex.homecatalogue.net/internal/validator"
)

func (app *application) createRecipeIngredientHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RecipeID     int64 `json:"recipe_id"`
		IngredientID int64 `json:"ingredient_id"`
		Amount       int32 `json:"amount"`
		Measurement  int64 `json:"measurement"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	recipeingredient := &data.RecipeIngredient{
		RecipeID:     input.RecipeID,
		IngredientID: input.IngredientID,
		Amount:       input.Amount,
		Measurement:  input.Measurement,
	}

	v := validator.New()

	if data.ValidateRecipeIngredient(v, recipeingredient); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.RecipeIngredients.Insert(recipeingredient)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/recipeingredients/%d/%d", recipeingredient.RecipeID, recipeingredient.IngredientID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"recipeingredients": recipeingredient}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRecipeIngredientHandler(w http.ResponseWriter, r *http.Request) {
	recipeid, ingredientid, err := app.readRecipeIngredientIDsParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recipeingredient, err := app.models.RecipeIngredients.Get(recipeid, ingredientid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipeingredient": recipeingredient}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRecipeIngredientHandler(w http.ResponseWriter, r *http.Request) {
	recipeid, ingredientid, err := app.readRecipeIngredientIDsParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recipeingredient, err := app.models.RecipeIngredients.Get(recipeid, ingredientid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		RecipeID     *int64 `json:"recipe_id"`
		IngredientID *int64 `json:"ingredient_id"`
		Amount       *int32 `json:"amount"`
		Measurement  *int64 `json:"measurement"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.RecipeID != nil {
		recipeingredient.RecipeID = *input.RecipeID
	}

	if input.IngredientID != nil {
		recipeingredient.IngredientID = *input.IngredientID
	}

	if input.Amount != nil {
		recipeingredient.Amount = *input.Amount
	}

	if input.Measurement != nil {
		recipeingredient.Measurement = *input.Measurement
	}

	v := validator.New()

	if data.ValidateRecipeIngredient(v, recipeingredient); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.RecipeIngredients.Update(recipeingredient)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipeingredient": recipeingredient}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRecipeIngredientHandler(w http.ResponseWriter, r *http.Request) {
	recipeid, ingredientid, err := app.readRecipeIngredientIDsParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.RecipeIngredients.Delete(recipeid, ingredientid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "recipe ingredient successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRecipeIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount      int
		Measurement int
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Amount = app.readInt(qs, "amount", 0, v)
	input.Measurement = app.readInt(qs, "measurement", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "recipe_id")

	input.Filters.SortSafelist = []string{"recipe_id", "ingredient_id", "amount", "measurement", "-recipe_id", "-ingredient_id", "-amount", "-measurement"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	recipeingredients, metadata, err := app.models.RecipeIngredients.GetAll(input.Amount, input.Measurement, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipeingredients": recipeingredients, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
