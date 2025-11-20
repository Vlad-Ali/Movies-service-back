package user

import (
	"context"

	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type Repository interface {
	GetByUserID(ctx context.Context, id object.UserID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
