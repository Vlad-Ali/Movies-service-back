package error

import "errors"

var (
	ErrReviewLikeIDCreatingIsNotValid = errors.New("review like ID is not valid")
	ErrReviewLikeIDAlreadyExists      = errors.New("review like ID already exists")
	ErrReviewLikeIsNotFound           = errors.New("review like is not found")
	ErrReviewLikeAlreadyExists        = errors.New("review like already exists")
)
