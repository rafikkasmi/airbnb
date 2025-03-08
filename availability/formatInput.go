package availability

import (
	"time"
)

// NewInputData creates a new InputData instance with current month and year
func NewInputData(roomId int64) InputData {
	now := time.Now()
	return InputData{
		RoomId:     roomId,
		StartMonth: int(now.Month()),
		StartYear:  now.Year(),
	}
}

// CoordinatesInput represents the coordinates for a map search
type CoordinatesInput struct {
	Ne CoordinatesValues
	Sw CoordinatesValues
}

// CoordinatesValues represents latitude and longitude values
type CoordinatesValues struct {
	Latitude float64
	Longitud float64
}
