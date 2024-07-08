package data

import (
	"time"
)

type AvailableItem struct {
	ID            int64       `json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	ExpirationAt  time.Time   `json:"expiration_at,omitempty"`
	LongName      string      `json:"long_name"`
	ShortName     string      `json:"short_name"`
	ItemType      ItemType    `json:"item_type"`
	Measurement   Measurement `json:"measurement"`
	ContainerSize int32       `json:"container_size"`
}
