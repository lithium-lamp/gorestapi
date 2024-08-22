package main

import (
	"errors"
	"fmt"
	"net/http"

	"householdingindex.homecatalogue.net/internal/data"
	"householdingindex.homecatalogue.net/internal/validator"
)

func (app *application) createKnownItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SerialNumber  int64    `json:"serial_number"`
		LongName      string   `json:"long_name"`
		ShortName     string   `json:"short_name"`
		Tags          []string `json:"tags"`
		ItemType      int64    `json:"item_type"`
		Measurement   int64    `json:"measurement"`
		ContainerSize int32    `json:"container_size"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	knownitem := &data.KnownItem{
		SerialNumber:  input.SerialNumber,
		LongName:      input.LongName,
		ShortName:     input.ShortName,
		Tags:          input.Tags,
		ItemType:      input.ItemType,
		Measurement:   input.Measurement,
		ContainerSize: input.ContainerSize,
	}

	v := validator.New()

	if data.ValidateKnownItem(v, knownitem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.KnownItems.Insert(knownitem)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/knownitems/%d", knownitem.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"knownitem": knownitem}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showKnownItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	knownitem, err := app.models.KnownItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"knownitem": knownitem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateKnownItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	knownitem, err := app.models.KnownItems.Get(id)
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
		SerialNumber  *int64   `json:"serial_number"`
		LongName      *string  `json:"long_name"`
		ShortName     *string  `json:"short_name"`
		Tags          []string `json:"tags"`
		ItemType      *int64   `json:"item_type"`
		Measurement   *int64   `json:"measurement"`
		ContainerSize *int32   `json:"container_size"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.SerialNumber != nil {
		knownitem.SerialNumber = *input.SerialNumber
	}

	if input.LongName != nil {
		knownitem.LongName = *input.LongName
	}

	if input.ShortName != nil {
		knownitem.ShortName = *input.ShortName
	}

	if input.Tags != nil {
		knownitem.Tags = input.Tags
	}

	if input.ItemType != nil {
		knownitem.ItemType = *input.ItemType
	}

	if input.Measurement != nil {
		knownitem.Measurement = *input.Measurement
	}

	if input.ContainerSize != nil {
		knownitem.ContainerSize = *input.ContainerSize
	}

	v := validator.New()

	if data.ValidateKnownItem(v, knownitem); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.KnownItems.Update(knownitem)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"knownitem": knownitem}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteKnownItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.KnownItems.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "known item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listKnownItemsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SerialNumber  int
		LongName      string
		ShortName     string
		Tags          []string
		ItemType      int
		Measurement   int
		ContainerSize int
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.SerialNumber = app.readInt(qs, "serial_number", 0, v)
	input.LongName = app.readString(qs, "long_name", "")
	input.ShortName = app.readString(qs, "short_name", "")
	input.Tags = app.readCSV(qs, "tags", []string{})
	input.ItemType = app.readInt(qs, "item_type", 0, v)
	input.Measurement = app.readInt(qs, "measurement", 0, v)
	input.ContainerSize = app.readInt(qs, "container_size", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "serial_number", "long_name", "short_name", "item_type", "measurement", "container_size", "-id", "-serial_number", "-long_name", "-short_name", "-item_type", "-measurement", "-container_size"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	knownitems, metadata, err := app.models.KnownItems.GetAll(input.SerialNumber, input.LongName, input.ShortName, input.Tags, input.ItemType, input.Measurement, input.ContainerSize, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"knownitems": knownitems, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
