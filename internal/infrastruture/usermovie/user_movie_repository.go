package usermovie

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	usermoviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/error"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/object"
	"github.com/lib/pq"
)

type UserMovieRepository struct {
	db *sql.DB
}

func NewUserMovieRepository(db *sql.DB) *UserMovieRepository {
	return &UserMovieRepository{db: db}
}

func (u *UserMovieRepository) Save(ctx context.Context, userMovie *usermoviedomain.UserMovie) error {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserMovieRepository.Save Error", "Error", err)
			return err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("UserMovieRepository.Save Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	var listType sql.NullString
	if userMovie.ListType() != usermoviedomain.ListTypeNone {
		listType = sql.NullString{
			String: string(userMovie.ListType()),
			Valid:  true,
		}
	}

	if userMovie.UserMovieID().IsEmpty() {
		query := `INSERT INTO user_movies (user_id, movie_id, list_type, user_rating) VALUES ($1, $2, $3, $4)
RETURNING id`
		var newID string
		err = tx.QueryRowContext(ctx, query, userMovie.UserID().ID(), userMovie.MovieID().ID(), listType, userMovie.UserRating()).Scan(&newID)
		if err != nil {
			slog.Error("UserMovieRepository.Save Error", "Error", err)
			return err
		}

		userMovieID, _ := object3.NewUserMovieID(newID)
		_ = userMovie.SetUserMovieID(userMovieID)
	} else {
		query := `
UPDATE user_movies SET list_type=$1, user_rating=$2 WHERE user_id=$3 AND movie_id=$4`
		result, execErr := tx.ExecContext(ctx, query, listType, userMovie.UserRating(), userMovie.UserID().ID(), userMovie.MovieID().ID())
		if execErr != nil {
			err = execErr
			slog.Error("UserMovieRepository.Save Error", "Error", execErr)
			return execErr
		}

		rowsAffected, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			slog.Error("UserMovieRepository.Save Error", "Error", rowsErr)
			err = rowsErr
			return rowsErr
		}

		if rowsAffected == 0 {
			return error2.ErrUserMovieIsNotFound
		}
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserMovieRepository.Save Error", "Error", commitErr)
			return commitErr
		}
	}
	return nil
}

func (u *UserMovieRepository) Delete(ctx context.Context, userMovie *usermoviedomain.UserMovie) error {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserMovieRepository.Delete Error", "Error", err)
			return err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("UserMovieRepository.Delete Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}
	query := `DELETE FROM user_movies WHERE user_id=$1 AND movie_id=$2`
	result, execErr := tx.ExecContext(ctx, query, userMovie.UserID().ID(), userMovie.MovieID().ID())
	if execErr != nil {
		err = execErr
		slog.Error("UserMovieRepository.Delete Error", "Error", err)
		return err
	}
	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		err = rowsErr
		slog.Error("UserMovieRepository.Delete Error", "Error", rowsErr)
		return err
	}
	if rowsAffected == 0 {
		return error2.ErrUserMovieIsNotFound
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserMovieRepository.Delete Error", "Error", commitErr)
			return commitErr
		}
	}

	return nil
}

func (u *UserMovieRepository) GetByUserAndMovie(ctx context.Context, userID object.UserID, movieID object2.MovieID) (*usermoviedomain.UserMovie, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserMovieRepository.GetByUserAndMovie Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("UserMovieRepository.GetByUserAndMovie Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}
	var nullListType sql.NullString
	userMovieModel := &UserMovieModel{}
	query := `SELECT id, user_id, movie_id, list_type, user_rating
FROM user_movies WHERE user_id=$1 AND movie_id=$2`
	err = tx.QueryRowContext(ctx, query, userID.ID(), movieID.ID()).Scan(&userMovieModel.ID, &userMovieModel.UserID, &userMovieModel.MovieID, &nullListType, &userMovieModel.UserRating)
	userMovieModel.ListType = nullListType.String
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error2.ErrUserMovieIsNotFound
		}
		slog.Error("UserMovieRepository.GetByUserAndMovie Error", "Error", err)
		return nil, err
	}
	userMovie, err := userMovieModel.ToDomain()
	if err != nil {
		slog.Error("UserMovieRepository.GetByUserAndMovie Error", "Error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserMovieRepository.GetByUserAndMovie Error", "Error", commitErr)
			return nil, commitErr
		}
	}

	return userMovie, nil
}

func (u *UserMovieRepository) GetMoviesByUserAndListType(ctx context.Context, userID object.UserID, listType usermoviedomain.ListType) ([]*usermoviedomain.MovieUserInfo, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserMovieRepository.GetMoviesByUserAndListType Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("UserMovieRepository.GetMoviesByUserAndListType Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	var movieUserInfos []*usermoviedomain.MovieUserInfo
	query := `SELECT 
            m.title,
            m.description,
            m.release_date,
            m.director,
            m.actors,
            m.genres,
            COALESCE((
                SELECT AVG(um2.user_rating)
                FROM user_movies um2 
                WHERE um2.movie_id = m.id AND um2.user_rating != 0
            ), 0) as rating,
            COALESCE(um.user_rating, 0) as user_rating
            FROM movies m
        LEFT JOIN user_movies um ON m.id = um.movie_id AND um.user_id = $1
        WHERE um.list_type IS NOT DISTINCT FROM $2`
	var nullListType sql.NullString
	if listType == usermoviedomain.ListTypeNone {
		nullListType = sql.NullString{Valid: false}
	} else {
		nullListType = sql.NullString{String: string(listType), Valid: true}
	}

	rows, err := tx.QueryContext(ctx, query, userID.ID(), nullListType)
	if err != nil {
		slog.Error("UserMovieRepository.GetMoviesByUserAndListType Error", "Error", err)
		return nil, err
	}
	defer rows.Close()
	movieUserInfos = make([]*usermoviedomain.MovieUserInfo, 0)
	for rows.Next() {
		movieUserInfo := &usermoviedomain.MovieUserInfo{Actors: make([]string, 0), Genres: make([]string, 0)}
		err = rows.Scan(&movieUserInfo.Title, &movieUserInfo.Description, &movieUserInfo.ReleaseDate, &movieUserInfo.Director,
			pq.Array(&movieUserInfo.Actors), pq.Array(&movieUserInfo.Genres), &movieUserInfo.Rating, &movieUserInfo.UserRating)
		if err != nil {
			slog.Error("UserMovieRepository.GetMoviesByUserAndListType Error", "Error", err)
			return nil, err
		}
		movieUserInfo.ListType = listType
		movieUserInfos = append(movieUserInfos, movieUserInfo)
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserMovieRepository.GetMoviesByUserAndListType Error", "Error", commitErr)
			return nil, commitErr
		}
	}
	return movieUserInfos, nil
}

func (u *UserMovieRepository) GetMovieByUserAndListType(ctx context.Context, userID object.UserID, movieID object2.MovieID, listType usermoviedomain.ListType) (*usermoviedomain.MovieUserInfo, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserMovieRepository.GetMovieByUserAndListType Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("UserMovieRepository.GetMovieByUserAndListType Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	movieUserInfo := &usermoviedomain.MovieUserInfo{Actors: make([]string, 0), Genres: make([]string, 0)}
	query := `SELECT 
            m.title,
            m.description,
            m.release_date,
            m.director,
            m.actors,
            m.genres,
            COALESCE((
                SELECT AVG(um2.user_rating)
                FROM user_movies um2 
                WHERE um2.movie_id = m.id AND um2.user_rating != 0
            ), 0) as rating,
            COALESCE(um.user_rating, 0) as user_rating
            FROM movies m
        LEFT JOIN user_movies um ON m.id = um.movie_id AND um.user_id = $1
        WHERE m.id = $2 AND um.list_type IS NOT DISTINCT FROM $3`
	var nullListType sql.NullString
	if listType == usermoviedomain.ListTypeNone {
		nullListType = sql.NullString{Valid: false}
	} else {
		nullListType = sql.NullString{String: string(listType), Valid: true}
	}

	err = tx.QueryRowContext(ctx, query, userID.ID(), movieID.ID(), nullListType).Scan(
		&movieUserInfo.Title, &movieUserInfo.Description, &movieUserInfo.ReleaseDate, &movieUserInfo.Director,
		pq.Array(&movieUserInfo.Actors), pq.Array(&movieUserInfo.Genres), &movieUserInfo.Rating, &movieUserInfo.UserRating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error2.ErrUserMovieIsNotFound
		}
		slog.Error("UserMovieRepository.GetMovieByUserAndListType Error", "Error", err)
		return nil, err
	}

	movieUserInfo.ListType = listType

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserMovieRepository.Delete Error", "Error", commitErr)
			return nil, commitErr
		}
	}
	return movieUserInfo, nil
}
