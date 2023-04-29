package model

// Metadata represents the movie metadata. Will be used by callers of our service
type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}
