package object

import (
	usermovieerror "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/error"
	"github.com/google/uuid"
)

type UserMovieID struct {
	id string
}

func NewUserMovieID(s string) (UserMovieID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return UserMovieID{}, usermovieerror.ErrIDCreatingIsNotValid
	}
	return UserMovieID{id: s}, nil
}

func (u UserMovieID) ID() string {
	return u.id
}

func (u UserMovieID) IsEmpty() bool {
	return u.id == ""
}
