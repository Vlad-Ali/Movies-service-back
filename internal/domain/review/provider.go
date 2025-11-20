package review

import (
	"context"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type Provider interface {
	ProvideMovieReviews(ctx context.Context, movieInfo object.MovieInfo) (string, error)
}
