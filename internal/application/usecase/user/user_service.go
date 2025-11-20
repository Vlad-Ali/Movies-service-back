package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/user/hasher"
	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/user/validation"
	userdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	usererror "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type UserService struct {
	tokenService   userdomain.TokenService
	userRepo       userdomain.Repository
	userTxManager  transactionmanager.TransactionManager[*userdomain.User]
	tokenTxManager transactionmanager.TransactionManager[string]
}

func NewUserService(tokenService userdomain.TokenService, userRepo userdomain.Repository, manager transactionmanager.TransactionManager[*userdomain.User], tokenTxManager transactionmanager.TransactionManager[string]) *UserService {
	return &UserService{tokenService: tokenService, userRepo: userRepo, userTxManager: manager, tokenTxManager: tokenTxManager}
}

func (u *UserService) GetUserByID(ctx context.Context, id object.UserID) (*userdomain.User, error) {
	return u.userTxManager.InTransaction(ctx, func(ctx context.Context) (*userdomain.User, error) {
		user, err := u.userRepo.GetByUserID(ctx, id)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to find user: %s, %v", id, err))
			return nil, err
		}
		slog.Debug("user found with id", "ID", id.ID())
		return user, nil
	})
}

func (u *UserService) Register(ctx context.Context, data object.UserRegistrationData) (*userdomain.User, error) {
	return u.userTxManager.InTransaction(ctx, func(ctx context.Context) (*userdomain.User, error) {
		err := uservalidation.ValidateUserRegistrationData(data)
		if err != nil {
			slog.Error("Validation failed", "error", err)
			return nil, err
		}

		ok, err := u.userRepo.ExistsByEmail(ctx, data.Email())
		if err != nil {
			slog.Error("Failed to check if email exists", "error", err)
			return nil, err
		}

		if ok {
			slog.Error("Email already exists")
			return nil, usererror.ErrUserEmailAlreadyExists
		}

		hashPassword, err := hasher.HashPassword(data.Password())
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			return nil, err
		}

		user := userdomain.NewUser(data.Username(), hashPassword, data.Email())
		user, err = u.userRepo.Save(ctx, user)
		if err != nil {
			slog.Error("Failed to create user", "error", err)
			return nil, usererror.ErrFailedToRegisterUser
		}
		slog.Debug("User created with id", "userID", user.ID().ID())
		return user, nil
	})
}

func (u *UserService) Authenticate(ctx context.Context, data object.AuthenticationData) (string, error) {
	return u.tokenTxManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		err := uservalidation.ValidateAuthenticationData(data)
		if err != nil {
			slog.Error("Validation failed", "error", err)
			return "", err
		}

		user, err := u.userRepo.GetByEmail(ctx, data.Email())
		if err != nil {
			slog.Error("failed to auth user error", "error", err)
			return "", err
		}

		ok := hasher.VerifyPassword(data.Password(), user.Password())
		if !ok {
			slog.Error("failed to auth user")
			return "", usererror.ErrInvalidPassword
		}

		token, err := u.tokenService.GenerateToken(ctx, user)
		if err != nil {
			slog.Error("failed to generate token error", "error", err)
			return "", err
		}
		slog.Debug("user is authenticated", "ID", user.ID().ID())
		return token, nil
	})

}
