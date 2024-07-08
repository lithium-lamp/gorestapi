package data

import (
	"fmt"
	"strconv"
)

type ItemType int32

func (it ItemType) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("this should be a string at index value %d", it)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
