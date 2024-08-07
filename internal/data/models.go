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
	Users          UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		AvailableItems: AvailableItemModel{DB: db},
		Users:          UserModel{DB: db},
	}
}
