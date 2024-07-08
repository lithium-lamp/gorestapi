package data

import (
	"fmt"
	"strconv"
)

type Measurement int32

func (m Measurement) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("this should be a string at index value %d", m)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
