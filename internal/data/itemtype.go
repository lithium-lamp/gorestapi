package data

/*

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidItemTypeFormat = errors.New("invalid item type format")

type ItemType int32

func (it ItemType) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", it)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (it *ItemType) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidItemTypeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidItemTypeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidItemTypeFormat
	}

	*it = ItemType(i)
	return nil
}
*/
