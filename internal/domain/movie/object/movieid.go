package object

import (
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/error"
	"github.com/google/uuid"
)

type MovieID struct {
	id string
}

func NewMovieID(s string) (MovieID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return MovieID{}, error2.ErrMovieIDCreatingIsNotValid
	}
	return MovieID{id: s}, nil
}

func (m MovieID) ID() string {
	return m.id
}

func (m MovieID) IsEmpty() bool {
	return m.id == ""
}
