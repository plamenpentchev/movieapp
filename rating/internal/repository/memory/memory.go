package memory

import (
	"context"

	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

// Repository defnes a rating repository
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new Repository
func New() *Repository {
	return &Repository{data: map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

// Get returns all ratings for a given record
func (r *Repository) Get(ctx context.Context, recordType model.RecordType, recordId model.RecordID) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrRecordTypeNotFound
	}
	if ratings, ok := r.data[recordType][recordId]; !ok || len(ratings) == 0 {
		return nil, repository.ErrRecordIdNotFound
	}
	return r.data[recordType][recordId], nil
}

// Put
func (r *Repository) Put(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}
	r.data[recordType][recordId] = append(r.data[recordType][recordId], *rating)
	return nil
}
