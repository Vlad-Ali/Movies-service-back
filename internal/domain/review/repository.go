package review

import (
	"context"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Repository interface {
	Save(ctx context.Context, review *Review) error
	Delete(ctx context.Context, review *Review) error
	GetReviewByUserAndMovie(ctx context.Context, userID object.UserID, movieID object2.MovieID) (*Review, error)
	GetReviewsByMovie(ctx context.Context, movieID object2.MovieID) ([]*ReviewInfo, error)
}
