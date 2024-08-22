package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Recipies          RecipeModel
	Ingredients       IngredientModel
	RecipeIngredients RecipeIngredientModel
	AvailableItems    AvailableItemModel
	KnownItems        KnownItemModel
	ItemTypes         ItemTypeModel
	Measurements      MeasurementModel
	Tags              TagModel
	Permissions       PermissionModel
	Tokens            TokenModel
	Users             UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Recipies:          RecipeModel{DB: db},
		Ingredients:       IngredientModel{DB: db},
		RecipeIngredients: RecipeIngredientModel{DB: db},
		AvailableItems:    AvailableItemModel{DB: db},
		KnownItems:        KnownItemModel{DB: db},
		ItemTypes:         ItemTypeModel{DB: db},
		Measurements:      MeasurementModel{DB: db},
		Tags:              TagModel{DB: db},
		Permissions:       PermissionModel{DB: db},
		Tokens:            TokenModel{DB: db},
		Users:             UserModel{DB: db},
	}
}
