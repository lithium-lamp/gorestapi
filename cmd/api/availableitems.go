package main

import (
	"fmt"
	"net/http"
	"time"

	"householdingindex.homecatalogue.net/internal/data"
	"householdingindex.homecatalogue.net/internal/validator"
)

func (app *application) createAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
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

	availableitem := data.AvailableItem{
		ID:            id,
		CreatedAt:     time.Now(),
		ExpirationAt:  time.Now().Add(time.Hour * 24 * 2),
		LongName:      "Isbergssallad Ca 440g Klass 1",
		ShortName:     "Isbergssallad",
		ItemType:      1,
		Measurement:   2,
		ContainerSize: 440,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"availableitem": availableitem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
