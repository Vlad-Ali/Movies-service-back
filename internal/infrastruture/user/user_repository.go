package user

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	userdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetByUserID(ctx context.Context, id object.UserID) (*userdomain.User, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserRepo.GetUserByUserID Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
			}
		}()
	}
	userModel := &UserModel{}
	query := "SELECT id, username, email, password_hash FROM users WHERE id = $1"
	err = tx.QueryRowContext(ctx, query, id.ID()).Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.Email,
		&userModel.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error2.ErrUserIsNotFound
		}
		slog.Error("UserRepo.GetUserByUserID Query Row Error", "Error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserRepo.GetUserByUserID Commit Error", "Error", commitErr)
			return nil, commitErr
		}
	}
	return userModel.ToDomain(), nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserRepo.GetByEmail Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
			}
		}()
	}
	userModel := &UserModel{}
	query := "SELECT id, username, email, password_hash FROM users WHERE email = $1"
	err = tx.QueryRowContext(ctx, query, email).Scan(
		&userModel.ID,
		&userModel.Username,
		&userModel.Email,
		&userModel.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error2.ErrUserIsNotFound
		}
		slog.Error("UserRepo.GetByEmail Query Row Error", "Error", err)
		return nil, err
	}
	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserRepo.GetByEmail Commit Error", "Error", commitErr)
			return nil, commitErr
		}
	}
	return userModel.ToDomain(), nil
}

func (u *UserRepository) Save(ctx context.Context, user *userdomain.User) (*userdomain.User, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserRepo.Save Begin Tx Error", "Error", err)
			return nil, err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
			}
		}()
	}

	if user.ID().IsEmpty() {
		query := `
INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)
RETURNING id`
		var newID string
		err = tx.QueryRowContext(ctx, query, user.Username(), user.Email(), user.Password()).Scan(&newID)
		if err != nil {
			slog.Error("UserRepo.Save Query Row Error", "Error", err)
			return nil, err
		}

		userID, _ := object.NewUserID(newID)
		_ = user.SetID(userID)
	} else {
		query := `
UPDATE users
SET username = $1, email = $2, password_hash = $3
WHERE id = $4`
		result, execErr := tx.ExecContext(ctx, query, user.Username(), user.Email(), user.Password(), user.ID())
		if execErr != nil {
			err = execErr
			slog.Error("UserRepo.Save Exec Error", "Error", err)
			return nil, err
		}

		rowsAffected, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			err = rowsErr
			slog.Error("UserRepo.Save RowsAffected Error", "Error", err)
			return nil, err
		}

		if rowsAffected == 0 {
			return nil, error2.ErrUserIsNotFound
		}
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserRepo.Save Commit Error", "Error", commitErr)
			return nil, commitErr
		}
	}
	return user, nil
}

func (u *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	var err error
	if !ok {
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("UserRepo.ExistsByEmail Begin Tx Error", "Error", err)
			return false, err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
			}
		}()
	}

	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	err = tx.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		slog.Error("UserRepo.ExistsByEmail Query Row Error", "Error", err)
		return false, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("UserRepo.ExistsByEmail Commit Error", "Error", commitErr)
			return false, commitErr
		}
	}
	return count > 0, nil
}
