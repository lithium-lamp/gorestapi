package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/availableitems", app.createAvailableItemHandler)
	router.HandlerFunc(http.MethodGet, "/v1/availableitems/:id", app.showAvailableItemHandler)
	router.HandlerFunc(http.MethodPut, "/v1/availableitems/:id", app.updateAvailableItemHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/availableitems/:id", app.deleteAvailableItemHandler)

	return router
}
