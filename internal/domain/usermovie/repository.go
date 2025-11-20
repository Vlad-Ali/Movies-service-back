package usermovie

import (
	"context"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Repository interface {
	Save(ctx context.Context, userMovie *UserMovie) error
	Delete(ctx context.Context, userMovie *UserMovie) error
	GetByUserAndMovie(ctx context.Context, userID object.UserID, movieID object2.MovieID) (*UserMovie, error)
	GetMoviesByUserAndListType(ctx context.Context, userID object.UserID, listType ListType) ([]*MovieUserInfo, error)
	GetMovieByUserAndListType(ctx context.Context, userID object.UserID, movieID object2.MovieID, listType ListType) (*MovieUserInfo, error)
}
