package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/recipies", app.requirePermission("recipies:read", app.listRecipiesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/recipies", app.requirePermission("recipies:write", app.createRecipeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/recipies/:id", app.requirePermission("recipies:read", app.showRecipeHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/recipies/:id", app.requirePermission("recipies:write", app.updateRecipeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/recipies/:id", app.requirePermission("recipies:write", app.deleteRecipeHandler))

	router.HandlerFunc(http.MethodGet, "/v1/ingredients", app.requirePermission("ingredients:read", app.listIngredientsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/ingredients", app.requirePermission("ingredients:write", app.createIngredientHandler))
	router.HandlerFunc(http.MethodGet, "/v1/ingredients/:id", app.requirePermission("ingredients:read", app.showIngredientHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/ingredients/:id", app.requirePermission("ingredients:write", app.updateIngredientHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/ingredients/:id", app.requirePermission("ingredients:write", app.deleteIngredientHandler))

	router.HandlerFunc(http.MethodGet, "/v1/recipeingredients", app.requirePermission("recipeingredients:read", app.listRecipeIngredientsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/recipeingredients", app.requirePermission("recipeingredients:write", app.createRecipeIngredientHandler))
	router.HandlerFunc(http.MethodGet, "/v1/recipeingredients/:recipe_id/:ingredient_id", app.requirePermission("recipeingredients:read", app.showRecipeIngredientHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/recipeingredients/:recipe_id/:ingredient_id", app.requirePermission("recipeingredients:write", app.updateRecipeIngredientHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/recipeingredients/:recipe_id/:ingredient_id", app.requirePermission("recipeingredients:write", app.deleteRecipeIngredientHandler))

	router.HandlerFunc(http.MethodGet, "/v1/availableitems", app.requirePermission("availableitems:read", app.listAvailableItemsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/availableitems", app.requirePermission("availableitems:write", app.createAvailableItemHandler))
	router.HandlerFunc(http.MethodGet, "/v1/availableitems/:id", app.requirePermission("availableitems:read", app.showAvailableItemHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/availableitems/:id", app.requirePermission("availableitems:write", app.updateAvailableItemHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/availableitems/:id", app.requirePermission("availableitems:write", app.deleteAvailableItemHandler))

	router.HandlerFunc(http.MethodGet, "/v1/knownitems", app.requirePermission("knownitems:read", app.listKnownItemsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/knownitems", app.requirePermission("knownitems:write", app.createKnownItemHandler))
	router.HandlerFunc(http.MethodGet, "/v1/knownitems/:id", app.requirePermission("knownitems:read", app.showKnownItemHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/knownitems/:id", app.requirePermission("knownitems:write", app.updateKnownItemHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/knownitems/:id", app.requirePermission("knownitems:write", app.deleteKnownItemHandler))

	router.HandlerFunc(http.MethodGet, "/v1/itemtypes", app.requirePermission("itemtypes:read", app.listItemTypesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/itemtypes", app.requirePermission("itemtypes:write", app.createItemTypeHandler))
	router.HandlerFunc(http.MethodGet, "/v1/itemtypes/:id", app.requirePermission("itemtypes:read", app.showItemTypeHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/itemtypes/:id", app.requirePermission("itemtypes:write", app.updateItemTypeHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/itemtypes/:id", app.requirePermission("itemtypes:write", app.deleteItemTypeHandler))

	router.HandlerFunc(http.MethodGet, "/v1/measurements", app.requirePermission("measurements:read", app.listMeasurementsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/measurements", app.requirePermission("measurements:write", app.createMeasurementHandler))
	router.HandlerFunc(http.MethodGet, "/v1/measurements/:id", app.requirePermission("measurements:read", app.showMeasurementHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/measurements/:id", app.requirePermission("measurements:write", app.updateMeasurementHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/measurements/:id", app.requirePermission("measurements:write", app.deleteMeasurementHandler))

	router.HandlerFunc(http.MethodGet, "/v1/tags", app.requirePermission("tags:read", app.listTagsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tags", app.requirePermission("tags:write", app.createTagHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tags/:id", app.requirePermission("tags:read", app.showTagHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/tags/:id", app.requirePermission("tags:write", app.updateTagHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tags/:id", app.requirePermission("tags:write", app.deleteTagHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/debug/vars", app.requirePermission("metrics:view", expvar.Handler().ServeHTTP))

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
