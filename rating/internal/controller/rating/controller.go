package rating

import (
	"context"
	"errors"
	"log"

	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

// ErrNotFound is returned when no ratings are found for a record.
var ErrNotFound = errors.New("ratings are not found for this record")

type ratingRepository interface {
	Get(ctx context.Context, recordType model.RecordType, recordId model.RecordID) ([]model.Rating, error)
	Put(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating *model.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepository
}

// New creates a new rating service controller.
func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if the are no records for it.
func (c *Controller) GetAggregatedRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordType, recordId)
	if err != nil && errors.Is(err, repository.ErrRecordIdNotFound) || errors.Is(err, repository.ErrRecordTypeNotFound) {
		log.Printf("Ratings not found due to an error: %v", err)
		return 0, ErrNotFound
	}
	sum := float64(0)
	if len(ratings) == 0 {
		log.Println("No ratings found")
		return sum, nil
	}

	for _, v := range ratings {
		sum += float64(v.RatingValue)
	}
	return sum / float64(len(ratings)), nil
}

func (c *Controller) PutRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating *model.Rating) error {
	return c.repo.Put(ctx, recordType, recordId, rating)
}
