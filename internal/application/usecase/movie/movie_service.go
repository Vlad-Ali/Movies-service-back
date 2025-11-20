package movie

import (
	"context"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	moviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type MovieService struct {
	moviesRepo          moviedomain.Repository
	movieTxManager      transactionmanager.TransactionManager[*moviedomain.Movie]
	movieSliceTxManager transactionmanager.TransactionManager[[]*moviedomain.Movie]
}

func NewMovieService(moviesRepo moviedomain.Repository, txManager transactionmanager.TransactionManager[*moviedomain.Movie], movieSliceTxManager transactionmanager.TransactionManager[[]*moviedomain.Movie]) *MovieService {
	return &MovieService{moviesRepo: moviesRepo, movieTxManager: txManager, movieSliceTxManager: movieSliceTxManager}
}

func (m *MovieService) FindByReleaseDateAndTitle(ctx context.Context, info object2.MovieInfo) (*moviedomain.Movie, error) {
	return m.movieTxManager.InTransaction(ctx, func(ctx context.Context) (*moviedomain.Movie, error) {
		movie, err := m.moviesRepo.GetByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("MovieService.FindByReleaseDateAndTitle failed to get movie", "error", err)
			return nil, err
		}
		slog.Debug("MovieService.FindByReleaseDateAndTitle movie successfully found by date and title")
		return movie, nil
	})
}

func (m *MovieService) GetAll(ctx context.Context) ([]*moviedomain.Movie, error) {
	return m.movieSliceTxManager.InTransaction(ctx, func(ctx context.Context) ([]*moviedomain.Movie, error) {
		movies, err := m.moviesRepo.GetAll(ctx)
		if err != nil {
			slog.Error("MovieService.GetAll failed to get movies", "error", err)
			return nil, err
		}
		slog.Debug("MovieService.GetAll all movies successfully found")
		return movies, nil
	})
}
