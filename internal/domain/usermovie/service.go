package usermovie

import (
	"context"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Service interface {
	SaveRating(ctx context.Context, userID object.UserID, info object2.MovieInfo, rating int) error
	SaveListType(ctx context.Context, userID object.UserID, info object2.MovieInfo, listType string) error
	FindMovieByUser(ctx context.Context, userID object.UserID, info object2.MovieInfo, listType string) (*MovieUserInfo, error)
	FindMoviesByUserAndListType(ctx context.Context, userID object.UserID, listType string) ([]*MovieUserInfo, error)
}
