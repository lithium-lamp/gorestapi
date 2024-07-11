package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidMeasurementFormat = errors.New("invalid measurement format")

type Measurement int32

func (m Measurement) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", m)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (m *Measurement) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidMeasurementFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidMeasurementFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidMeasurementFormat
	}

	*m = Measurement(i)
	return nil
}
