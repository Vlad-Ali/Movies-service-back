package user

import (
	"context"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Service interface {
	GetUserByID(ctx context.Context, id object.UserID) (*User, error)
	Register(ctx context.Context, data object.UserRegistrationData) (*User, error)
	Authenticate(ctx context.Context, data object.AuthenticationData) (*object.AuthResponse, error)
}
