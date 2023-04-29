package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/pkg/model"
)

// Handler defines a rating service controller
type Handler struct {
	ctrl *rating.Controller
}

// New creates a new rating service controller
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// Handle handles PUT and GET /rating requests
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordId := model.RecordID(r.FormValue("id"))
	if recordId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(r.Context(), recordType, recordId)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err = json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("response encode error: %v", err)
		}
	case http.MethodPost:
		userId := model.UserID(r.FormValue("userId"))
		if userId == "" {
			log.Println("wrong user id value. Empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			log.Printf("wrong value. Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rating := &model.Rating{
			RecordID:    "",
			RecordType:  "",
			UserID:      userId,
			RatingValue: model.RatingValue(v),
		}

		err = h.ctrl.PutRating(r.Context(), recordType, recordId, rating)
		if err != nil {
			log.Printf("repository put error. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
