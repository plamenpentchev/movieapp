package model

//RecordID defines a record id. Together with RecordType identifies unnique record
//across all types.
type RecordID string

//RecordType defines a record type. Combined with RecordID identifies unnique records
//across all types
type RecordType string

//Existing record types
const (
	RecordTypeMovie = RecordType("movie")
)

//UserID defines a user id.
type UserID string

//RatingValue defines a value of a rating record.
type RatingValue float64

//Rating defines an individual rating created by the user
type Rating struct {
	RecordID    RecordID    `json:"recordId"`
	RecordType  RecordType  `json:"recordType"`
	UserID      UserID      `json:"userId"`
	RatingValue RatingValue `json:"value"`
}
