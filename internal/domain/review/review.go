package review

import (
	"time"

	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
)

var (
	emptyTextLen = 0
	maxTextLen   = 200
)

func ValidateReviewText(text string) error {
	if !(len(text) < maxTextLen && len(text) > emptyTextLen) {
		return error2.ErrReviewTextValidationError
	}
	return nil
}

type Review struct {
	id          object.ReviewID
	userID      object2.UserID
	movieID     object3.MovieID
	text        string
	writingDate time.Time
	userRating  int
}

func NewReview(userID object2.UserID, movieID object3.MovieID) *Review {
	return &Review{userID: userID, movieID: movieID}
}

func (r *Review) ID() object.ReviewID {
	return r.id
}

func (r *Review) SetID(id object.ReviewID) error {
	if r.ID().IsEmpty() {
		r.id = id
		return nil
	}
	return error2.ErrReviewIDAlreadyExists
}

func (r *Review) UserID() object2.UserID {
	return r.userID
}

func (r *Review) MovieID() object3.MovieID {
	return r.movieID
}

func (r *Review) Text() string {
	return r.text
}

func (r *Review) SetText(text string) error {
	err := ValidateReviewText(text)
	if err != nil {
		return err
	}
	r.text = text
	return nil
}

func (r *Review) WritingDate() time.Time {
	return r.writingDate
}

func (r *Review) SetWritingDate(date time.Time) {
	r.writingDate = date
}

func (r *Review) UserRating() int {
	return r.userRating
}

func (r *Review) SetUserRating(rating int) error {
	err := usermovie.ValidateUserRating(rating)
	if err != nil {
		return err
	}
	r.userRating = rating
	return nil
}
