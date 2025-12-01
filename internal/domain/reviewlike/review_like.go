package reviewlike

import (
	object3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike/object"
	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type ReviewLike struct {
	id       object.ReviewLikeID
	userID   object2.UserID
	reviewID object3.ReviewID
}

func NewReviewLike(userID object2.UserID, reviewID object3.ReviewID) *ReviewLike {
	return &ReviewLike{userID: userID, reviewID: reviewID}
}

func (like *ReviewLike) ID() object.ReviewLikeID {
	return like.id
}

func (like *ReviewLike) UserID() object2.UserID {
	return like.userID
}

func (like *ReviewLike) ReviewID() object3.ReviewID {
	return like.reviewID
}

func (like *ReviewLike) SetID(id object.ReviewLikeID) error {
	if !like.id.IsEmpty() {
		return error2.ErrReviewLikeIDAlreadyExists
	}
	like.id = id
	return nil
}
