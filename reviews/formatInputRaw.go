package reviews

type variables struct {
	Id                string            `json:"id"`
	PdpReviewsRequest pdpReviewsRequest `json:"pdpReviewsRequest"`
}

type pdpReviewsRequest struct {
	FieldSelector     string `json:"fieldSelector"`
	ForPreview        bool   `json:"forPreview"`
	Limit             int    `json:"limit"`
	Offset            string `json:"offset"`
	SortingPreference string `json:"sortingPreference"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}
