package error

import "errors"

var (
	ErrInvalidRating            = errors.New("rating must be between 0 and 10")
	ErrIDCreatingIsNotValid     = errors.New("id is not valid")
	ErrUserMovieIsNotFound      = errors.New("user movie is not found")
	ErrListTypeIsIncorrect      = errors.New("list type is incorrect")
	ErrUserMovieIDAlreadyExists = errors.New("user movie ID already exists")
)
