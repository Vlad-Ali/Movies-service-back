package review

import (
	"context"
	"time"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Service interface {
	SaveReview(ctx context.Context, userID object.UserID, movieInfo object2.MovieInfo, text string, writingDate time.Time) error
	DeleteReview(ctx context.Context, userID object.UserID, info object2.MovieInfo) error
	GetUserReview(ctx context.Context, userID object.UserID, info object2.MovieInfo) (*Review, error)
	GetReviewsByMovie(ctx context.Context, info object2.MovieInfo) ([]*ReviewInfo, error)
	GetReviewsByMovieForUser(ctx context.Context, info object2.MovieInfo, userID object.UserID) ([]*ReviewInfo, error)
}
