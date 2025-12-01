package reviewlike

import (
	"context"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Service interface {
	LikeReview(ctx context.Context, userID object.UserID, reviewID object2.ReviewID) error
	UnLikeReview(ctx context.Context, userID object.UserID, reviewID object2.ReviewID) error
}
