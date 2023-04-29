package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"movieexample.com/movie/internal/gateway"
	"movieexample.com/rating/pkg/model"
)

// Gateway defines an HTTP gateway for a rating service.
type Gateway struct {
	addr string
}

// New creates a new gateway to a rating service.
func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (model.RatingValue, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ratinng", g.addr), nil)
	if err != nil {
		return model.RatingValue(0.0), err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("id", fmt.Sprintf("%s", recordId))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.RatingValue(0.0), err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return model.RatingValue(0.0), gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return model.RatingValue(0.0), fmt.Errorf("non-2xx response: %v", resp)
	}
	var v model.RatingValue

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return model.RatingValue(0.0), err
	}
	return v, nil
}

// PutRating writes a new rating
func (g *Gateway) PutRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/rating", g.addr), nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", string(recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.RatingValue))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
