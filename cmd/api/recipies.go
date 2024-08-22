package main

import (
	"errors"
	"fmt"
	"net/http"

	"householdingindex.homecatalogue.net/internal/data"
	"householdingindex.homecatalogue.net/internal/validator"
)

func (app *application) createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		CookingSteps    []string `json:"cooking_steps"`
		CookTimeMinutes int32    `json:"cook_time_minutes"`
		Portions        int32    `json:"portions"`
		Tags            []string `json:"tags"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	recipe := &data.Recipe{
		Name:            input.Name,
		Description:     input.Description,
		CookingSteps:    input.CookingSteps,
		CookTimeMinutes: input.CookTimeMinutes,
		Portions:        input.Portions,
		Tags:            input.Tags,
	}

	v := validator.New()

	if data.ValidateRecipe(v, recipe); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Recipies.Insert(recipe)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/recipies/%d", recipe.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"recipe": recipe}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recipe, err := app.models.Recipies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recipe, err := app.models.Recipies.Get(id)
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
		Name            *string  `json:"name"`
		Description     *string  `json:"description"`
		CookingSteps    []string `json:"cooking_steps"`
		CookTimeMinutes *int32   `json:"cook_time_minutes"`
		Portions        *int32   `json:"portions"`
		Tags            []string `json:"tags"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		recipe.Name = *input.Name
	}

	if input.Description != nil {
		recipe.Description = *input.Description
	}

	if input.CookingSteps != nil {
		recipe.CookingSteps = input.CookingSteps
	}

	if input.CookTimeMinutes != nil {
		recipe.CookTimeMinutes = *input.CookTimeMinutes
	}

	if input.Portions != nil {
		recipe.Portions = *input.Portions
	}

	if input.Tags != nil {
		recipe.Tags = input.Tags
	}

	v := validator.New()

	if data.ValidateRecipe(v, recipe); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Recipies.Update(recipe)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipe": recipe}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Recipies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "recipies successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRecipiesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name            string
		Description     string
		CookingSteps    []string
		CookTimeMinutes int
		Portions        int
		Tags            []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Description = app.readString(qs, "description", "")
	input.CookingSteps = app.readCSV(qs, "cooking_steps", []string{})
	input.CookTimeMinutes = app.readInt(qs, "cook_time_minutes", 0, v)
	input.Portions = app.readInt(qs, "portions", 0, v)
	input.Tags = app.readCSV(qs, "tags", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "description", "cook_time_minutes", "portions", "-id", "-name", "-description", "-cook_time_minutes", "-portions"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	recipies, metadata, err := app.models.Recipies.GetAll(input.Name, input.Description, input.CookingSteps, input.CookTimeMinutes, input.Portions, input.Tags, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"recipies": recipies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
