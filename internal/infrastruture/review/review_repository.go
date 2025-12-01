package reviewrepo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/error"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Save(ctx context.Context, review *reviewdomain.Review) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewRepo.Save Begin Tx Error", "Error", err)
			return err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewRepo.Save Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	if review.ID().IsEmpty() {
		var newID string
		query := `INSERT INTO reviews (user_id, movie_id, text, writing_date) VALUES 
                                                                ($1, $2, $3, $4)
                                                                RETURNING id`
		execErr := tx.QueryRowContext(ctx, query, review.UserID().ID(), review.MovieID().ID(), review.Text(), review.WritingDate()).Scan(&newID)
		if execErr != nil {
			slog.Error("ReviewRepo.Save Exec Error", "Error", execErr, "UserID", review.UserID().ID(), "MovieID", review.MovieID().ID())
			err = execErr
			return err
		}

		reviewID, _ := object3.NewReviewID(newID)
		_ = review.SetID(reviewID)
	} else {
		query := `UPDATE reviews SET text = $1, writing_date = $2 WHERE id = $3`

		result, execErr := tx.ExecContext(ctx, query, review.Text(), review.WritingDate(), review.ID().ID())
		if execErr != nil {
			slog.Error("ReviewRepo.Save Exec Error", "Error", execErr)
			err = execErr
			return err
		}

		rowsAffected, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			slog.Error("ReviewRepo.Save RowsAffected Error", "Error", rowsErr)
			err = rowsErr
			return err
		}

		if rowsAffected == 0 {
			return error2.ErrReviewNotFound
		}
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewRepo.Save Commit Error", "Error", commitErr)
			_ = tx.Rollback()
			return commitErr
		}
	}

	return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, review *reviewdomain.Review) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewRepo.Delete Begin Tx Error", "Error", err)
			return err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewRepo.Delete Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	query := `DELETE FROM reviews WHERE user_id = $1 AND movie_id = $2`
	result, execErr := tx.ExecContext(ctx, query, review.UserID().ID(), review.MovieID().ID())
	if execErr != nil {
		slog.Error("ReviewRepo.Delete Exec Error", "Error", execErr)
		err = execErr
		return err
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		slog.Error("ReviewRepo.Delete RowsAffected Error", "Error", rowsErr)
		err = rowsErr
		return err
	}

	if rowsAffected == 0 {
		return error2.ErrReviewNotFound
	}
	return nil
}

func (r *ReviewRepository) GetReviewByUserAndMovie(ctx context.Context, userID object.UserID, movieID object2.MovieID) (*reviewdomain.Review, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewRepo.GetReviewByUserAndMovie Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewRepo.GetReviewByUserAndMovie Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	reviewModel := &ReviewModel{}
	query := `SELECT id, user_id, movie_id, text, writing_date FROM reviews WHERE user_id = $1 AND movie_id = $2`
	err = tx.QueryRowContext(ctx, query, userID.ID(), movieID.ID()).Scan(&reviewModel.ID, &reviewModel.UserID, &reviewModel.MovieID, &reviewModel.Text, &reviewModel.WritingDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, error2.ErrReviewNotFound
	} else if err != nil {
		slog.Error("ReviewRepo.GetReviewByUserAndMovie Query Error", "Error", err)
		return nil, err
	}
	review, err := reviewModel.ToDomain()
	if err != nil {
		slog.Error("ReviewRepo.GetReviewByUserAndMovie ToDomain Error", "Error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewRepo.GetReviewByUserAndMovie Commit Error", "Error", err)
			_ = tx.Rollback()
			return nil, commitErr
		}
	}

	return review, nil
}

func (r *ReviewRepository) GetReviewsByMovie(ctx context.Context, movieID object2.MovieID) ([]*reviewdomain.ReviewInfo, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewRepo.GetReviewsByMovie Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewRepo.GetReviewsByMovie Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	query := `SELECT id, (SELECT u.username FROM users AS u WHERE u.id = r.user_id), r.text, r.writing_date, COALESCE((SELECT um.user_rating FROM user_movies AS um
              WHERE um.user_id = r.user_id AND um.movie_id = r.movie_id), 0), (SELECT COUNT(*) FROM review_likes AS rl WHERE rl.review_id = r.id) as likes FROM reviews AS r
              WHERE r.movie_id = $1
              ORDER BY likes DESC 
              LIMIT 100`

	reviews := make([]*reviewdomain.ReviewInfo, 0)
	rows, err := tx.QueryContext(ctx, query, movieID.ID())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("ReviewRepo.GetReviewsByMovie Query Error", "Error", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		reviewInfo := &reviewdomain.ReviewInfo{}
		var date time.Time
		err = rows.Scan(&reviewInfo.ID, &reviewInfo.Username, &reviewInfo.Text, &date, &reviewInfo.UserRating, &reviewInfo.Likes)
		reviewInfo.ReviewYear = date.Year()
		reviewInfo.ReviewMonth = int(date.Month())
		reviewInfo.ReviewDay = date.Day()
		if err != nil {
			slog.Error("ReviewRepo.GetReviewsByMovie Scan", "Error", err)
			return nil, err
		}

		reviews = append(reviews, reviewInfo)
	}

	err = rows.Err()
	if err != nil {
		slog.Error("ReviewRepo.GetReviewsByMovie Rows Error", "Error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewRepo.GetReviewsByMovie Commit Error", "Error", err)
			_ = tx.Rollback()
			return nil, commitErr
		}
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewByMovieForUser(ctx context.Context, movieID object2.MovieID, userID object.UserID) ([]*reviewdomain.ReviewInfo, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewRepo.GetReviewByMovieForUser Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewRepo.GetReviewByMovieForUser Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	query := `SELECT id, (SELECT u.username FROM users AS u WHERE u.id = r.user_id), r.text, r.writing_date, COALESCE((SELECT um.user_rating FROM user_movies AS um
              WHERE um.user_id = r.user_id AND um.movie_id = r.movie_id), 0), EXISTS(SELECT 1 FROM review_likes AS rl WHERE rl.review_id = r.id AND rl.user_id = $2),  (SELECT COUNT(*) FROM review_likes AS rl WHERE rl.review_id = r.id) as likes FROM reviews AS r
              WHERE r.movie_id = $1
              ORDER BY likes DESC 
              LIMIT 100`

	reviews := make([]*reviewdomain.ReviewInfo, 0)
	rows, err := tx.QueryContext(ctx, query, movieID.ID(), userID.ID())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("ReviewRepo.GetReviewByMovieForUser Query Error", "Error", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		reviewInfo := &reviewdomain.ReviewInfo{}
		var date time.Time
		err = rows.Scan(&reviewInfo.ID, &reviewInfo.Username, &reviewInfo.Text, &date, &reviewInfo.UserRating, &reviewInfo.IsLiked, &reviewInfo.Likes)
		reviewInfo.ReviewYear = date.Year()
		reviewInfo.ReviewMonth = int(date.Month())
		reviewInfo.ReviewDay = date.Day()
		if err != nil {
			slog.Error("ReviewRepo.GetReviewByMovieForUser Scan", "Error", err)
			return nil, err
		}

		reviews = append(reviews, reviewInfo)
	}

	err = rows.Err()
	if err != nil {
		slog.Error("ReviewRepo.GetReviewByMovieForUser Rows Error", "Error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewRepo.GetReviewByMovieForUser Commit Error", "Error", err)
			_ = tx.Rollback()
			return nil, commitErr
		}
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewByID(ctx context.Context, reviewID object3.ReviewID) (*reviewdomain.Review, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("ReviewRepo.GetReviewByID Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					slog.Error("ReviewRepo.GetReviewByID Rollback Error", "Error", rollbackErr)
				}
			}
		}()
	}

	reviewModel := &ReviewModel{}
	query := `SELECT id, user_id, movie_id, text, writing_date FROM reviews WHERE id = $1`
	err = tx.QueryRowContext(ctx, query, reviewID.ID()).Scan(&reviewModel.ID, &reviewModel.UserID, &reviewModel.MovieID, &reviewModel.Text, &reviewModel.WritingDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, error2.ErrReviewNotFound
	} else if err != nil {
		slog.Error("ReviewRepo.GetReviewByID Query Error", "Error", err)
		return nil, err
	}
	review, err := reviewModel.ToDomain()
	if err != nil {
		slog.Error("ReviewRepo.GetReviewByID ToDomain Error", "Error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("ReviewRepo.GetReviewByID Commit Error", "Error", err)
			_ = tx.Rollback()
			return nil, commitErr
		}
	}

	return review, nil
}
