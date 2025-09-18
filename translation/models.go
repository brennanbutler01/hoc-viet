package translation

type Output struct {
	Body any
}

type Match struct {
	ID             any     `json:"id"`
	Segment        string  `json:"segment"`
	Translation    string  `json:"translation"`
	Source         string  `json:"source"`
	Target         string  `json:"target"`
	Quality        any     `json:"quality"`
	UsageCount     int     `json:"usage-count"`
	CreatedBy      string  `json:"created-by"`
	LastUpdatedBy  string  `json:"last-updated-by"`
	CreateDate     string  `json:"create-date"`
	LastUpdateDate string  `json:"last-update-date"`
	Match          float64 `json:"match"`
}

type ResponseData struct {
	TranslatedText string  `json:"translatedText"`
	Match          float64 `json:"match"`
}

// MyMemoryResponse is the complete response from MyMemory API
type MyMemoryResponse struct {
	ResponseData   ResponseData `json:"responseData"`
	QuotaFinished  bool         `json:"quotaFinished"`
	ResponseStatus int          `json:"responseStatus"`
	Matches        []Match      `json:"matches"`
}
