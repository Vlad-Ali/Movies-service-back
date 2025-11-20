package movie

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	moviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/lib/pq"
)

type MovieRepository struct {
	db *sql.DB
}

func (m *MovieRepository) GetIDByReleaseDateAndTitle(ctx context.Context, title string, year int, month int, day int) (object.MovieID, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("MovieRepo.IDByReleaseDateAndTitle Begin Tx Error", "Error", err)
			return object.MovieID{}, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("MovieRepo.IDByReleaseDateAndTitle Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	query := `SELECT id FROM movies
WHERE title = $1 AND release_date = $2`
	var id string
	err = tx.QueryRowContext(ctx, query, title, date).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return object.MovieID{}, error2.ErrMovieIsNotFound
	} else if err != nil {
		slog.Error("MovieRepo.GetIDByReleaseDateAndTitle row Scan", "Error", err)
		return object.MovieID{}, err
	}

	movieID, _ := object.NewMovieID(id)

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("MovieRepo.GetIDByReleaseDateAndTitle Commit Error", "Error", err)
			return object.MovieID{}, commitErr
		}
	}

	return movieID, nil
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (m *MovieRepository) GetAll(ctx context.Context) ([]*moviedomain.Movie, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("MovieRepo.GetAll Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("MovieRepo.GetAll Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	query := `SELECT m.id, m.title, m.description, m.release_date, m.director, m.actors, m.genres,
       COALESCE((SELECT AVG(um.user_rating)
                 FROM user_movies as um
                 WHERE um.movie_id = m.id AND um.user_rating !=0 ), 0)
                 FROM movies as m`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		slog.Error("MovieRepo.GetAll Query Error", "Error", err)
		return nil, err
	}
	defer rows.Close()

	moviesInfos := make([]*moviedomain.Movie, 0)
	for rows.Next() {
		var id string
		movie := &moviedomain.Movie{Actors: make([]string, 0), Genres: make([]string, 0)}
		err = rows.Scan(&id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Director, pq.Array(&movie.Actors), pq.Array(&movie.Genres), &movie.Rating)
		if err != nil {
			slog.Error("MovieRepo.GetAll Scan Error", "Error", err)
			return nil, err
		}
		movieID, _ := object.NewMovieID(id)
		_ = movie.SetID(movieID)
		moviesInfos = append(moviesInfos, movie)
	}
	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("MovieRepo.GetAll Commit Error", "Error", commitErr)
			return nil, commitErr
		}
	}
	return moviesInfos, nil
}

func (m *MovieRepository) GetByReleaseDateAndTitle(ctx context.Context, title string, year int, month int, day int) (*moviedomain.Movie, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("MovieRepo.GetByReleaseDateAndTitle Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("MovieRepo.GetByReleaseDateAndTitle Rollback Error", "Error", err)
				}
			}
		}()
	}

	var id string
	releaseDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	movie := &moviedomain.Movie{Actors: make([]string, 0), Genres: make([]string, 0)}
	query := `
SELECT m.id, m.title, m.description, m.release_date, m.director, m.actors, m.genres,
       COALESCE((SELECT AVG(um.user_rating)
                 FROM user_movies as um
                 WHERE um.movie_id = m.id AND um.user_rating !=0 ), 0)
                 FROM movies as m
WHERE m.title = $1 AND m.release_date = $2
`
	err = tx.QueryRowContext(ctx, query, title, releaseDate).Scan(&id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Director, pq.Array(&movie.Actors), pq.Array(&movie.Genres), &movie.Rating)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("MovieRepo.GetByReleaseDateAndTitle Error", "Error", err, "Title", title, "Date", releaseDate)
		return nil, error2.ErrMovieIsNotFound
	}
	if err != nil {
		slog.Error("MovieRepo.GetByReleaseDateAndTitle Error row Scan", "Error", err)
		return nil, err
	}

	movieID, _ := object.NewMovieID(id)
	_ = movie.SetID(movieID)
	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("MovieRepo.GetByReleaseDateAndTitle Commit Error", "Error", err)
			return nil, commitErr
		}
	}
	return movie, nil
}
