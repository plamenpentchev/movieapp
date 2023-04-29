package model

import (
	metadataModel "movieexample.com/metadata/pkg/model"
	ratingModel "movieexample.com/rating/pkg/model"
)

// MovieDetails includes movie metadata informatio, its aggregated rating.
type MovieDetails struct {
	Metadata metadataModel.Metadata  `json:"metadata"`
	Rating   ratingModel.RatingValue `json:"rating,omitempty"`
}
