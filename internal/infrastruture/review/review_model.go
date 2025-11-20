package reviewrepo

import (
	"time"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type ReviewModel struct {
	ID          string
	UserID      string
	MovieID     string
	Text        string
	WritingDate time.Time
	UserRating  int
}

func (r *ReviewModel) ToDomain() (*reviewdomain.Review, error) {
	userID, err := object.NewUserID(r.UserID)
	if err != nil {
		return nil, err
	}

	movieID, err := object2.NewMovieID(r.MovieID)
	if err != nil {
		return nil, err
	}

	review := reviewdomain.NewReview(userID, movieID)
	reviewID, err := object3.NewReviewID(r.ID)
	if err != nil {
		return nil, err
	}

	err = review.SetID(reviewID)
	if err != nil {
		return nil, err
	}

	err = review.SetText(r.Text)
	if err != nil {
		return nil, err
	}

	review.SetWritingDate(r.WritingDate)
	err = review.SetUserRating(r.UserRating)
	if err != nil {
		return nil, err
	}
	return review, nil
}
