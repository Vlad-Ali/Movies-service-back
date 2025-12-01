package reviewservice

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	moviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type ReviewService struct {
	movieRepo        moviedomain.Repository
	reviewRepo       reviewdomain.Repository
	txUser           transactionmanager.TransactionUser
	reviewTxManager  transactionmanager.TransactionManager[*reviewdomain.Review]
	reviewsTxManager transactionmanager.TransactionManager[[]*reviewdomain.ReviewInfo]
}

func NewReviewService(movieRepo moviedomain.Repository, reviewRepo reviewdomain.Repository, txUser transactionmanager.TransactionUser, reviewTxManager transactionmanager.TransactionManager[*reviewdomain.Review], reviewsTxManager transactionmanager.TransactionManager[[]*reviewdomain.ReviewInfo]) *ReviewService {
	return &ReviewService{movieRepo: movieRepo, reviewRepo: reviewRepo, txUser: txUser, reviewTxManager: reviewTxManager, reviewsTxManager: reviewsTxManager}
}

func (r *ReviewService) SaveReview(ctx context.Context, userID object.UserID, movieInfo object2.MovieInfo, text string, writingDate time.Time) error {
	return r.txUser.UseTransaction(ctx, func(ctx context.Context) error {
		err := reviewdomain.ValidateReviewText(text)
		if err != nil {
			slog.Error("ReviewSrv.SaveReview Error review validation failed", "error", err)
			return err
		}

		movieID, err := r.movieRepo.GetIDByReleaseDateAndTitle(ctx, movieInfo.Title, movieInfo.Year, movieInfo.Month, movieInfo.Day)
		if err != nil {
			slog.Error("ReviewSrv.SaveReview Error while getting movie", "error", err)
			return err
		}

		review, err := r.reviewRepo.GetReviewByUserAndMovie(ctx, userID, movieID)
		if err != nil && !errors.Is(err, error2.ErrReviewNotFound) {
			slog.Error("ReviewSrv.SaveReview Error while getting review", "error", err)
			return err
		} else if errors.Is(err, error2.ErrReviewNotFound) {
			slog.Debug("ReviewSrv.SaveReview Error", "error", err, "movieID", movieID.ID(), "userID", userID.ID())
			review = reviewdomain.NewReview(userID, movieID)
		}

		err = review.SetText(text)
		if err != nil {
			slog.Error("ReviewSrv.SaveReview Error while setting text", "error", err)
			return err
		}

		review.SetWritingDate(writingDate)
		err = r.reviewRepo.Save(ctx, review)
		if err != nil {
			slog.Error("ReviewSrv.SaveReview Error while saving review", "error", err)
			return err
		}
		return nil
	})
}

func (r *ReviewService) DeleteReview(ctx context.Context, userID object.UserID, info object2.MovieInfo) error {
	return r.txUser.UseTransaction(ctx, func(ctx context.Context) error {
		movieID, err := r.movieRepo.GetIDByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("ReviewSrv.DeleteReview Error while getting movie", "error", err)
			return err
		}

		review := reviewdomain.NewReview(userID, movieID)
		slog.Debug("ReviewSrv.DeleteReview", "movieID", movieID, "userID", userID.ID())
		err = r.reviewRepo.Delete(ctx, review)
		if err != nil {
			slog.Error("ReviewSrv.DeleteReview Error while deleting review", "error", err)
			return err
		}

		return nil
	})
}

func (r *ReviewService) GetUserReview(ctx context.Context, userID object.UserID, info object2.MovieInfo) (*reviewdomain.Review, error) {
	return r.reviewTxManager.InTransaction(ctx, func(ctx context.Context) (*reviewdomain.Review, error) {
		movieID, err := r.movieRepo.GetIDByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("ReviewSrv.GetUserReview Error while getting movie", "error", err)
			return nil, err
		}

		review, err := r.reviewRepo.GetReviewByUserAndMovie(ctx, userID, movieID)
		if err != nil && !errors.Is(err, error2.ErrReviewNotFound) {
			slog.Error("ReviewSrv.GetUserReview Error while getting review", "error", err)
			return nil, err
		} else if errors.Is(err, error2.ErrReviewNotFound) {
			return &reviewdomain.Review{}, nil
		}

		return review, nil
	})
}

func (r *ReviewService) GetReviewsByMovie(ctx context.Context, info object2.MovieInfo) ([]*reviewdomain.ReviewInfo, error) {
	return r.reviewsTxManager.InTransaction(ctx, func(ctx context.Context) ([]*reviewdomain.ReviewInfo, error) {
		movieID, err := r.movieRepo.GetIDByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("ReviewSrv.GetUserReviews Error while getting movie", "error", err)
			return nil, err
		}

		reviews, err := r.reviewRepo.GetReviewsByMovie(ctx, movieID)
		if err != nil {
			slog.Error("ReviewSrv.GetUserReviews Error while getting reviews", "error", err)
			return nil, err
		}

		return reviews, nil
	})
}

func (r *ReviewService) GetReviewsByMovieForUser(ctx context.Context, info object2.MovieInfo, userID object.UserID) ([]*reviewdomain.ReviewInfo, error) {
	return r.reviewsTxManager.InTransaction(ctx, func(ctx context.Context) ([]*reviewdomain.ReviewInfo, error) {
		movieID, err := r.movieRepo.GetIDByReleaseDateAndTitle(ctx, info.Title, info.Year, info.Month, info.Day)
		if err != nil {
			slog.Error("ReviewSrv.GetUserReviews Error while getting movie", "error", err)
			return nil, err
		}

		reviews, err := r.reviewRepo.GetReviewByMovieForUser(ctx, movieID, userID)
		if err != nil {
			slog.Error("ReviewSrv.GetUserReviews Error while getting reviews", "error", err)
			return nil, err
		}

		return reviews, nil
	})
}
