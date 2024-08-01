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
		ExpirationAt  time.Time        `json:"expiration_at"`
		LongName      string           `json:"long_name"`
		ShortName     string           `json:"short_name"`
		ItemType      data.ItemType    `json:"item_type"`
		Measurement   data.Measurement `json:"measurement"`
		ContainerSize int32            `json:"container_size"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	availableitem := &data.AvailableItem{
		ExpirationAt:  input.ExpirationAt,
		LongName:      input.LongName,
		ShortName:     input.ShortName,
		ItemType:      input.ItemType,
		Measurement:   input.Measurement,
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
		ExpirationAt  time.Time        `json:"expiration_at"`
		LongName      string           `json:"long_name"`
		ShortName     string           `json:"short_name"`
		ItemType      data.ItemType    `json:"item_type"`
		Measurement   data.Measurement `json:"measurement"`
		ContainerSize int32            `json:"container_size"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	availableitem.ExpirationAt = input.ExpirationAt
	availableitem.LongName = input.LongName
	availableitem.ShortName = input.ShortName
	availableitem.ItemType = input.ItemType
	availableitem.Measurement = input.Measurement
	availableitem.ContainerSize = input.ContainerSize

	v := validator.New()

	if data.ValidateAvailableItem(v, availableitem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.AvailableItems.Update(availableitem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
