package object

import (
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike/error"
	"github.com/google/uuid"
)

type ReviewLikeID struct {
	id string
}

func NewReviewLikeID(id string) (ReviewLikeID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ReviewLikeID{}, error2.ErrReviewLikeIDCreatingIsNotValid
	}
	return ReviewLikeID{id}, nil
}

func (obj ReviewLikeID) ID() string {
	return obj.id
}

func (obj ReviewLikeID) IsEmpty() bool {
	return obj.id == ""
}
