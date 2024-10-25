package types

// APIConfig configures additional traits about this API.
type APIConfig struct {
	BaseURL string `json:"baseURL"` // Base URL of this API. All endpoints are relative to this endpoint
}
