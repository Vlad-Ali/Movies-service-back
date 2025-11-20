package error

import "errors"

var (
	ErrMovieIDCreatingIsNotValid = errors.New("movie id is not valid")
	ErrMovieIsNotFound           = errors.New("movie not found")
	ErrMovieDataValidationFailed = errors.New("movie data validation failed")
	ErrMovieIDAlreadyExists      = errors.New("movie id already exists")
)
