package reviewlike

import (
	"context"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	reviewlikedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type ReviewLikeService struct {
	reviewRepository     reviewdomain.Repository
	reviewLikeRepository reviewlikedomain.Repository
	txUser               transactionmanager.TransactionUser
}

func NewReviewLikeService(reviewRepository reviewdomain.Repository, reviewLikeRepository reviewlikedomain.Repository, txUser transactionmanager.TransactionUser) *ReviewLikeService {
	return &ReviewLikeService{reviewRepository: reviewRepository, reviewLikeRepository: reviewLikeRepository, txUser: txUser}
}

func (r *ReviewLikeService) LikeReview(ctx context.Context, userID object.UserID, reviewID object3.ReviewID) error {
	return r.txUser.UseTransaction(ctx, func(ctx context.Context) error {
		_, err := r.reviewRepository.GetReviewByID(ctx, reviewID)
		if err != nil {
			slog.Error("ReviewLikeService.LikeReview Get review error", "error", err)
			return err
		}

		exists, err := r.reviewLikeRepository.Exists(ctx, userID, reviewID)
		if err != nil {
			slog.Error("ReviewLikeService.LikeReview Exists error", "error", err)
			return err
		}

		if exists {
			slog.Error("ReviewLikeService.LikeReview like already exists", "userID", userID, "reviewID", reviewID)
			return error2.ErrReviewLikeAlreadyExists
		}

		err = r.reviewLikeRepository.Like(ctx, userID, reviewID)
		if err != nil {
			slog.Error("ReviewLikeService.LikeReview like error", "error", err)
			return err
		}

		return nil
	})
}

func (r *ReviewLikeService) UnLikeReview(ctx context.Context, userID object.UserID, reviewID object3.ReviewID) error {
	return r.txUser.UseTransaction(ctx, func(ctx context.Context) error {
		_, err := r.reviewRepository.GetReviewByID(ctx, reviewID)
		if err != nil {
			slog.Error("ReviewLikeService.UnLikeReview Get review error", "error", err)
			return err
		}

		exists, err := r.reviewLikeRepository.Exists(ctx, userID, reviewID)
		if err != nil {
			slog.Error("ReviewLikeService.UnLikeReview Exists error", "error", err)
			return err
		}

		if !exists {
			slog.Error("ReviewLikeService.UnLikeReview like does not exist", "userID", userID, "reviewID", reviewID)
			return error2.ErrReviewLikeIsNotFound
		}

		err = r.reviewLikeRepository.UnLike(ctx, userID, reviewID)
		if err != nil {
			slog.Error("ReviewLikeService.UnLikeReview unlike error", "error", err)
			return err
		}

		return nil
	})
}
