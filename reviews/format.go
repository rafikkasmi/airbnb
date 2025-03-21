package reviews

type Data struct {
	RoomID           int64
	Badges           []string
	Name             string
	Title            string
	Type             string
	Kind             string
	Category         string
	Rating           Rating
	Coordinates      Coordinates
	Fee              Fee
	Price            PriceData
	LongStayDiscount Price
	Images           []Img
	Reviews          []Review
}

type Review struct {
	ID                 string `csv:"ID"`
	Comments           string `csv:"Comments"`
	TranslatedComments string `csv:"TranslatedComments"`
	Language           string `csv:"Language"`
	CreatedAt          string `csv:"CreatedAt"`
	LocalizedDate      string `csv:"LocalizedDate"`
	Rating             int    `csv:"Rating"`
	Reviewer           User   `csv:"Reviewer"`
	Host               User   `csv:"Host"`
	Highlight          string `csv:"Highlight"`
}

type User struct {
	ID            string
	FirstName     string
	FullName      string
	IsSuperhost   bool
	ProfilePicURL string
	ProfilePath   string
}
type Fee struct {
	Cleaning Price
	Airbn    Price
}
type PriceData struct {
	Total     Price
	Unit      UnitPrice
	BreakDown []PriceBreakDown
}
type PriceBreakDown struct {
	Description    string
	Amount         float32
	CurrencySymbol string
}
type Rating struct {
	Value       float32
	ReviewCount int
}
type Coordinates struct {
	Latitude float64
	Longitud float64
}

type Img struct {
	URL         string
	ContentType string
	Extension   string
	Content     []byte `json:"-"`
}
type UnitPrice struct {
	Qualifier string
	Discount  float32
	Price
}
type Price struct {
	Amount         float32
	CurrencySymbol string
}
