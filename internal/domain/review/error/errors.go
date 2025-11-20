package error

import "errors"

var (
	ErrReviewNotFound            = errors.New("review not found")
	ErrReviewIDAlreadyExists     = errors.New("review id already exists")
	ErrReviewTextValidationError = errors.New("review text validation error")
)
