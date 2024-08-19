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
	AvailableItems AvailableItemModel
	KnownItems     KnownItemModel
	ItemTypes      ItemTypeModel
	Measurements   MeasurementModel
	Permissions    PermissionModel
	Tokens         TokenModel
	Users          UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		AvailableItems: AvailableItemModel{DB: db},
		KnownItems:     KnownItemModel{DB: db},
		ItemTypes:      ItemTypeModel{DB: db},
		Measurements:   MeasurementModel{DB: db},
		Permissions:    PermissionModel{DB: db},
		Tokens:         TokenModel{DB: db},
		Users:          UserModel{DB: db},
	}
}
