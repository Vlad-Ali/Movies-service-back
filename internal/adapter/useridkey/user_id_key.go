package useridkey

import (
	"log/slog"
	"net/http"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type UserIDKey struct{}

func ExtractUserIdFromReq(r *http.Request) (object.UserID, error) {
	id := r.Context().Value(UserIDKey{}).(string)
	userID, err := object.NewUserID(id)
	if err != nil {
		slog.Error("Error while extracting user id from request", "error", err)
		return object.UserID{}, err
	}
	return userID, nil
}
