package details

type Data struct {
	RoomID               int64            `csv:"RoomID"`
	Title                string           `csv:"Title"`
	URL                  string           `csv:"URL"`
	RoomType             string           `csv:"RoomType"`
	Language             string           `csv:"Language"`
	HomeTier             int              `csv:"HomeTier"`
	PersonCapacity       int              `csv:"PersonCapacity"`
	IsSuperHost          bool             `csv:"IsSuperHost"`
	Price                Price            `csv:"Price"`
	Rating               Rating           `csv:"Rating"`
	Coordinates          Coordinates      `csv:"Coordinates"`
	Host                 Host             `csv:"Host"`
	CoHosts              []Cohost         `csv:"CoHosts"`
	SubDescription       SubDescription   `csv:"SubDescription"`
	Description          string           `csv:"Description"`
	Highlights           []Highlight      `csv:"Highlights"`
	Amenities            []AmenityGroup   `csv:"Amenities"`
	HouseRules           HouseRules       `csv:"HouseRules"`
	LocationDescriptions []LocationDetail `csv:"LocationDescriptions"`
	Images               []Img            `csv:"Images"`

	// // Availability and night-related fields
	// MinimumNights        int    `json:"minimum_nights"`
	// MaximumNights        int    `json:"maximum_nights"`
	// MinimumMinimumNights int    `json:"minimum_minimum_nights"`
	// MaximumMinimumNights int    `json:"maximum_minimum_nights"`
	// MinimumMaximumNights int    `json:"minimum_maximum_nights"`
	// MaximumMaximumNights int    `json:"maximum_maximum_nights"`
	// MinimumNightsAvgNtm  int    `json:"minimum_nights_avg_ntm"`
	// MaximumNightsAvgNtm  int    `json:"maximum_nights_avg_ntm"`
	// CalendarUpdated      string `json:"calendar_updated"`
	// HasAvailability      bool   `json:"has_availability"`
	// Availability30       int    `json:"availability_30"`
	// Availability60       int    `json:"availability_60"`
	// Availability90       int    `json:"availability_90"`
}
type Price struct {
	Amount         float32
	CurrencySymbol string
	Qualifier      string
}
type Cohost struct {
	ID   string
	Name string
}
type Host struct {
	ID          string
	Name        string
	JoinedOn    string
	Description string
}
type HouseRules struct {
	Aditional string
	General   []HouseRule
}
type Img struct {
	Title       string
	URL         string
	ContentType string
	Extension   string
	Content     []byte `json:"-"`
}
type HouseRule struct {
	Title  string
	Values []HouseRuleValue
}
type HouseRuleValue struct {
	Title string
	Icon  string
}
type LocationDetail struct {
	Title   string
	Content string
}
type Rating struct {
	Accuracy          float32
	Checking          float32
	CleaningLiness    float32
	Comunication      float32
	Location          float32
	Value             float32
	GuestSatisfaction float32
	ReviewCount       int
}
type Coordinates struct {
	Latitude float64
	Longitud float64
}
type SubDescription struct {
	Title string
	Items []string
}
type AmenityGroup struct {
	Title  string
	Values []Amenity
}
type Amenity struct {
	Title     string
	Subtitle  string
	Available bool
	Icon      string
}
type Highlight struct {
	Title    string
	Subtitle string
	Icon     string
}
