package reviewlike

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type ReviewLikeRepository struct {
	db *sql.DB
}

func NewReviewLikeRepository(db *sql.DB) *ReviewLikeRepository {
	return &ReviewLikeRepository{db: db}
}

func (r ReviewLikeRepository) Like(ctx context.Context, userID object.UserID, reviewID object3.ReviewID) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewLikeRepo.Like Begin Tx Error", "Error", err)
			return err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewLikeRepo.Like Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}
	query := `INSERT INTO review_likes (user_id, review_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, query, userID.ID(), reviewID.ID())
	if err != nil {
		slog.Error("ReviewLikeRepo.Like Exec Error", "Error", err)
		return err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewLikeRepo.Like Commit Error", "Error", commitErr)
			_ = tx.Rollback()
			return commitErr
		}
	}

	return nil
}

func (r ReviewLikeRepository) UnLike(ctx context.Context, userID object.UserID, reviewID object3.ReviewID) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewLikeRepo.UnLike Begin Tx Error", "Error", err)
			return err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewLikeRepo.UnLike Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}
	query := `DELETE FROM review_likes WHERE user_id = $1 AND review_id = $2`
	_, err = tx.ExecContext(ctx, query, userID.ID(), reviewID.ID())
	if err != nil {
		slog.Error("ReviewLikeRepo.UnLike Exec Error", "Error", err)
		return err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewLikeRepo.UnLike Commit Error", "Error", commitErr)
			_ = tx.Rollback()
			return commitErr
		}
	}

	return nil
}

func (r ReviewLikeRepository) Exists(ctx context.Context, userID object.UserID, reviewID object3.ReviewID) (bool, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewLikeRepo.Exists Begin Tx Error", "Error", err)
			return false, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewLikeRepo.Exists Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM review_likes WHERE user_id = $1 AND review_id = $2)`
	err = tx.QueryRowContext(ctx, query, userID.ID(), reviewID.ID()).Scan(&exists)
	if err != nil {
		slog.Error("ReviewLikeRepo.Exists Exec Error", "Error", err)
		return false, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewLikeRepo.Exists Commit Error", "Error", commitErr)
			_ = tx.Rollback()
			return false, commitErr
		}
	}

	return exists, nil
}
