package model

type ShortenRequest struct {
	URL string `json:"url"`
}

type BatchShortenRequestInstance struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchShortenRequest []BatchShortenRequestInstance

type ShortenResponse struct {
	Result string `json:"result"`
}

type BatchShortenResponseInstance struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchShortenResponse []BatchShortenResponseInstance

type UserURLsResponceInstance struct {
	ShortURL    string `json:"short_url"`
	OriginalURL URL    `json:"original_url"`
}

type UserURLsResponce []UserURLsResponceInstance
