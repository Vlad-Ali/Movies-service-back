package movie

import (
	"context"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type Service interface {
	FindByReleaseDateAndTitle(ctx context.Context, info object.MovieInfo) (*Movie, error)
	GetAll(ctx context.Context) ([]*Movie, error)
}
