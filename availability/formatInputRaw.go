package availability

type variables struct {
	Request request `json:"request"`
}

type request struct {
	Count     int    `json:"count"`
	ListingId string `json:"listingId"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}
