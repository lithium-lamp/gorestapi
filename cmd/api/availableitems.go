package main

import (
	"fmt"
	"net/http"
	"time"

	"householdingindex.homecatalogue.net/internal/data"
)

func (app *application) createAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new available item")
}

func (app *application) showAvailableItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
