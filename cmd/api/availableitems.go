package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"householdingindex.homecatalogue.net/internal/data"
	"householdingindex.homecatalogue.net/internal/validator"
)

func (app *application) createAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		KnownItemsID  int64     `json:"knownitems_id"`
		ExpirationAt  time.Time `json:"expiration_at"`
		ContainerSize int32     `json:"container_size"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	availableitem := &data.AvailableItem{
		KnownItemsID:  input.KnownItemsID,
		ExpirationAt:  input.ExpirationAt,
		ContainerSize: input.ContainerSize,
	}

	v := validator.New()

	if data.ValidateAvailableItem(v, availableitem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.AvailableItems.Insert(availableitem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/availableitems/%d", availableitem.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"availableitem": availableitem}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	availableitem, err := app.models.AvailableItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"availableitem": availableitem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	availableitem, err := app.models.AvailableItems.Get(id)
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
		KnownItemsID  *int64     `json:"knownitems_id"`
		ExpirationAt  *time.Time `json:"expiration_at"`
		ContainerSize *int32     `json:"container_size"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.KnownItemsID != nil {
		availableitem.KnownItemsID = *input.KnownItemsID
	}

	if input.ExpirationAt != nil {
		availableitem.ExpirationAt = *input.ExpirationAt
	}

	if input.ContainerSize != nil {
		availableitem.ContainerSize = *input.ContainerSize
	}

	v := validator.New()

	if data.ValidateAvailableItem(v, availableitem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.AvailableItems.Update(availableitem)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"availableitem": availableitem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.AvailableItems.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "available item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAvailableItemsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		KnownItemsID  int
		ExpirationAt  time.Time
		ContainerSize int
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.KnownItemsID = app.readInt(qs, "knownitems_id", 0, v)

	input.ExpirationAt = app.readTime(qs, "expiration_at", time.Time{}, v)

	input.ContainerSize = app.readInt(qs, "container_size", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "knownitems_id", "expiration_at", "container_size", "-id", "-knownitems_id", "-expiration_at", "-container_size"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	availableitems, metadata, err := app.models.AvailableItems.GetAll(input.KnownItemsID, input.ExpirationAt, input.ContainerSize, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"availableitems": availableitems, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
