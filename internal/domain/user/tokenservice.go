package user

import (
	"context"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type TokenService interface {
	GenerateToken(ctx context.Context, user *User) (string, error)
	ValidateToken(ctx context.Context, token string) (object.UserID, error)
}
