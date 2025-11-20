package usermovie

import (
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	usermoviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/object"
)

type UserMovieModel struct {
	ID         string
	UserID     string
	MovieID    string
	ListType   string
	UserRating int
}

func (u *UserMovieModel) ToDomain() (*usermoviedomain.UserMovie, error) {
	userID, err := object.NewUserID(u.UserID)
	if err != nil {
		return nil, err
	}
	movieID, err := object2.NewMovieID(u.MovieID)
	if err != nil {
		return nil, err
	}
	userMovieID, err := object3.NewUserMovieID(u.ID)
	if err != nil {
		return nil, err
	}
	userMovie := usermoviedomain.NewUserMovie(userID, movieID)
	err = userMovie.SetUserMovieID(userMovieID)
	if err != nil {
		return nil, err
	}
	err = userMovie.SetRating(u.UserRating)
	if err != nil {
		return nil, err
	}
	listType, err := usermoviedomain.ValidateAndGetListType(u.ListType)
	if err != nil {
		return nil, err
	}
	userMovie.SetListType(listType)
	return userMovie, nil
}
