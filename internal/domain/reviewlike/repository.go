package reviewlike

import (
	"context"

	object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Repository interface {
	Like(ctx context.Context, userID object.UserID, reviewID object2.ReviewID) error
	UnLike(ctx context.Context, userID object.UserID, reviewID object2.ReviewID) error
	Exists(ctx context.Context, userID object.UserID, reviewID object2.ReviewID) (bool, error)
}
