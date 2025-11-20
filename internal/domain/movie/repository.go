package movie

import (
	"context"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*Movie, error)
	GetByReleaseDateAndTitle(ctx context.Context, title string, year int, month int, day int) (*Movie, error)
	GetIDByReleaseDateAndTitle(ctx context.Context, title string, year int, month int, day int) (object.MovieID, error)
}
