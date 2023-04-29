package repository

import "errors"

//ErrRecordTypeNotFound is returned when no record was found given the record type.
var ErrRecordTypeNotFound = errors.New("record type not found")

//ErrRecordIdNotFound is returned when no record was found given the record id.
var ErrRecordIdNotFound = errors.New("record id not found")
