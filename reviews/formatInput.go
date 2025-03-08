package reviews

type InputData struct {
	RoomId int64
}

type CoordinatesInput struct {
	Ne CoordinatesValues
	Sw CoordinatesValues
}
type CoordinatesValues struct {
	Latitude float64
	Longitud float64
}
