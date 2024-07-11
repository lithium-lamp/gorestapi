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

	v := validator.New()

	v.Check(input.LongName != "", "long_name", "must be provided")
	v.Check(len(input.LongName) <= 500, "long_name", "must not be more than 500 bytes long")
	//v.Check(validator.Unique(input.LongName), "long_name", "must not contain duplicate values")

	v.Check(input.ShortName != "", "short_name", "must be provided")
	v.Check(len(input.ShortName) <= 100, "short_name", "must not be more than 100 bytes long")

	v.Check(input.ItemType != 0, "item_type", "must be provided")
	v.Check(input.ItemType >= 1, "item_type", "must be greater than 0")
	v.Check(input.ItemType <= 6, "year", "must not be greater than 6") //TEMP VALUE

	v.Check(input.Measurement != 0, "measurement", "must be provided")
	v.Check(input.Measurement >= 1, "measurement", "must be greater than 0")
	v.Check(input.Measurement <= 6, "measurement", "must not be greater than 6")

	v.Check(input.ContainerSize >= 0, "container_size", "must be at least 0")
	v.Check(input.ContainerSize <= 100000, "container_size", "must not be more than 100000 units")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
