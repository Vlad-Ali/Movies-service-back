package object

import (
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/error"
	"github.com/google/uuid"
)

type ReviewID struct {
	id string
}

func NewReviewID(id string) (ReviewID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ReviewID{}, error2.ErrReviewIDAlreadyExists
	}
	return ReviewID{id: id}, nil
}

func (r ReviewID) ID() string {
	return r.id
}

func (r ReviewID) IsEmpty() bool {
	return r.id == ""
}
