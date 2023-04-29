package movie

import (
	"context"
	"errors"

	metadatModel "movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	movieModel "movieexample.com/movie/pkg/model"
	ratingModel "movieexample.com/rating/pkg/model"
)

// ErrNotFound
var ErrNotFound = errors.New("movie data not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordId ratingModel.RecordID, recordType ratingModel.RecordType) (ratingModel.RatingValue, error)
	PutRating(ctx context.Context, recordId ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error
}

type metadatGateway interface {
	Get(ctx context.Context, id string) (*metadatModel.Metadata, error)
}

// Controller defines a movie service controller
type Controller struct {
	ratingGateway  ratingGateway
	metadatGateway metadatGateway
}

// New creates a new movie service controller.
func New(ratingGateway ratingGateway, metadatGateway metadatGateway) *Controller {
	return &Controller{ratingGateway, metadatGateway}

}

// Get returns the movie details including the aggregate rating and movie metadata.
func (c *Controller) Get(ctx context.Context, id string) (*movieModel.MovieDetails, error) {
	metadata, err := c.metadatGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &movieModel.MovieDetails{Metadata: *metadata}

	ratings, err := c.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		// just proceed in this case, it's ok to not have ratings yet.
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = ratings
	}
	return details, nil
}
