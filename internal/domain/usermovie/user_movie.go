package usermovie

import (
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/error"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/object"
)

type ListType string

const (
	ListTypeFavorite  ListType = "favorite"
	ListTypeWatchlist ListType = "watchlist"
	ListTypeNone      ListType = ""
)

const EmptyRating = 0

func ValidateAndGetListType(listType string) (ListType, error) {
	switch listType {
	case string(ListTypeFavorite):
		return ListTypeFavorite, nil
	case string(ListTypeWatchlist):

		return ListTypeWatchlist, nil
	case string(ListTypeNone):

		return ListTypeNone, nil
	default:

		return ListTypeNone, error2.ErrListTypeIsIncorrect
	}
}

func ValidateUserRating(rating int) error {
	if rating < 0 || rating > 10 {
		return error2.ErrInvalidRating
	}
	return nil
}

type UserMovie struct {
	id         object3.UserMovieID
	userID     object.UserID
	movieID    object2.MovieID
	listType   ListType
	userRating int
}

func NewUserMovie(userID object.UserID, movieID object2.MovieID) *UserMovie {
	return &UserMovie{
		userID:     userID,
		movieID:    movieID,
		listType:   ListTypeNone,
		userRating: 0,
	}
}

func (u *UserMovie) SetRating(rating int) error {
	err := ValidateUserRating(rating)
	if err != nil {
		return err
	}
	u.userRating = rating
	return nil
}

func (u *UserMovie) SetListType(listType ListType) {
	u.listType = listType
}

func (u *UserMovie) ListType() ListType {
	return u.listType
}

func (u *UserMovie) UserRating() int {
	return u.userRating
}

func (um *UserMovie) IsFavorite() bool {
	return um.listType == ListTypeFavorite
}

func (um *UserMovie) IsInWatchlist() bool {
	return um.listType == ListTypeWatchlist
}

func (um *UserMovie) HasRating() bool {
	return um.userRating > 0
}

func (um *UserMovie) UserMovieID() object3.UserMovieID {
	return um.id
}

func (um *UserMovie) SetUserMovieID(id object3.UserMovieID) error {

	if um.id.IsEmpty() {
		um.id = id
		return nil
	}
	return error2.ErrUserMovieIDAlreadyExists
}

func (um *UserMovie) MovieID() object2.MovieID {
	return um.movieID
}

func (um *UserMovie) UserID() object.UserID {
	return um.userID
}

func (um *UserMovie) IsEmpty() bool {
	return um.listType == ListTypeNone && um.userRating == EmptyRating
}
