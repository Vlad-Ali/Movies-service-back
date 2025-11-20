package usermovie

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	moviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	usermoviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/error"
)

type UserMovieService struct {
	movieInfosTxManager transactionmanager.TransactionManager[[]*usermoviedomain.MovieUserInfo]
	movieInfoTxManager  transactionmanager.TransactionManager[*usermoviedomain.MovieUserInfo]
	txUser              transactionmanager.TransactionUser
	moviesRepo          moviedomain.Repository
	userMovieRepo       usermoviedomain.Repository
}

func NewUserMovieService(moviesRepo moviedomain.Repository, userMovieRepo usermoviedomain.Repository, movieInfosTxManager *transactionmanager.TransactionManagerImpl[[]*usermoviedomain.MovieUserInfo], movieInfoTxManager transactionmanager.TransactionManager[*usermoviedomain.MovieUserInfo], txUser transactionmanager.TransactionUser) *UserMovieService {
	return &UserMovieService{
		moviesRepo:          moviesRepo,
		userMovieRepo:       userMovieRepo,
		movieInfoTxManager:  movieInfoTxManager,
		movieInfosTxManager: movieInfosTxManager,
		txUser:              txUser,
	}
}

func (u *UserMovieService) SaveRating(ctx context.Context, userID object.UserID, info object2.MovieInfo, rating int) error {
	return u.txUser.UseTransaction(ctx, func(ctx context.Context) error {
		movie, err := u.moviesRepo.GetByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("UMSvc.SaveRating GetByReleaseDateAndTitle failed", "error", err)
			return err
		}

		userMovie, err := u.userMovieRepo.GetByUserAndMovie(ctx, userID, movie.ID())
		if err != nil && !errors.Is(err, error2.ErrUserMovieIsNotFound) {
			slog.Error("UMSvc.SaveRating GetByUserAndMovie failed", "error", err)
			return err
		} else if errors.Is(err, error2.ErrUserMovieIsNotFound) {
			userMovie = usermoviedomain.NewUserMovie(userID, movie.ID())
		}
		err = userMovie.SetRating(rating)
		if err != nil {
			slog.Error("UMSvc.SaveRating SetRating failed", "error", err)
			return err
		}

		if !userMovie.UserMovieID().IsEmpty() && userMovie.IsEmpty() {
			err = u.userMovieRepo.Delete(ctx, userMovie)
			if err != nil {
				slog.Error("UMSvc.SaveRating Delete failed", "error", err)
				return err
			}
		} else {
			err = u.userMovieRepo.Save(ctx, userMovie)
			if err != nil {
				slog.Error("UMSvc.SaveRating SaveUserMovie failed", "error", err)
				return err
			}
		}
		slog.Debug("UMSvc.SaveRating user movie rating saved")
		return nil
	})
}

func (u *UserMovieService) SaveListType(ctx context.Context, userID object.UserID, info object2.MovieInfo, listType string) error {
	movieListType, err := usermoviedomain.ValidateAndGetListType(listType)
	if err != nil {
		slog.Error("UMSvc.SaveListType Validation failed", "error", err)
		return err
	}
	return u.txUser.UseTransaction(ctx, func(ctx context.Context) error {
		movie, err := u.moviesRepo.GetByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("UMSvc.SaveListType GetByReleaseDateAndTitle failed", "error", err)
			return err
		}

		userMovie, err := u.userMovieRepo.GetByUserAndMovie(ctx, userID, movie.ID())
		if err != nil && !errors.Is(err, error2.ErrUserMovieIsNotFound) {
			slog.Error("UMSvc.SaveListType GetByUserAndMovie failed", "error", err)
			return err
		} else if errors.Is(err, error2.ErrUserMovieIsNotFound) {
			userMovie = usermoviedomain.NewUserMovie(userID, movie.ID())
		}

		userMovie.SetListType(movieListType)
		if !userMovie.UserMovieID().IsEmpty() && userMovie.IsEmpty() {
			err = u.userMovieRepo.Delete(ctx, userMovie)
			if err != nil {
				slog.Error("UMSvc.SaveListType DeleteUserMovie failed", "error", err)
				return err
			}
		} else {
			err = u.userMovieRepo.Save(ctx, userMovie)
			if err != nil {
				slog.Error("UMSvc.SaveListType SaveUserMovie failed", "error", err)
				return err
			}
		}
		slog.Debug("UMSvc.SaveListType user movie list type saved")
		return nil
	})
}

func (u *UserMovieService) FindMovieByUser(ctx context.Context, userID object.UserID, info object2.MovieInfo, listType string) (*usermoviedomain.MovieUserInfo, error) {
	movieListType, err := usermoviedomain.ValidateAndGetListType(listType)
	if err != nil {
		slog.Error("UMSvc.FindMovieByUser Validation failed", "error", err)
		return nil, err
	}
	return u.movieInfoTxManager.InTransaction(ctx, func(ctx context.Context) (*usermoviedomain.MovieUserInfo, error) {
		movie, err := u.moviesRepo.GetByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)

		if err != nil {
			slog.Error("UMSvc.FindMovieByUser GetByReleaseDateAndTitle failed", "error", err)
			return nil, err
		}

		movieUserInfo, err := u.userMovieRepo.GetMovieByUserAndListType(ctx, userID, movie.ID(), movieListType)
		if err != nil {
			slog.Error("UMSvc.FindMovieByUser Failed to get MovieInfo by user", "error", err)
			return nil, err
		}
		slog.Debug("UMSvc.FindMovieByUser successfully found  movie by user")
		return movieUserInfo, nil
	})
}

func (u *UserMovieService) FindMoviesByUserAndListType(ctx context.Context, userID object.UserID, listType string) ([]*usermoviedomain.MovieUserInfo, error) {
	movieListType, err := usermoviedomain.ValidateAndGetListType(listType)
	if err != nil {
		slog.Error("UMSvc.FindMoviesByUserAndListType Validation failed", "error", err)
		return nil, err
	}
	return u.movieInfosTxManager.InTransaction(ctx, func(ctx context.Context) ([]*usermoviedomain.MovieUserInfo, error) {
		movieUserInfos, err := u.userMovieRepo.GetMoviesByUserAndListType(ctx, userID, movieListType)
		if err != nil {
			slog.Error("UMSvc.FindMoviesByUserAndListType GetMoviesByUserAndListType failed", "error", err)
			return nil, err
		}
		slog.Debug("UMSvc.FindMoviesByUserAndListType movies successfully found by user")
		return movieUserInfos, nil
	})
}
