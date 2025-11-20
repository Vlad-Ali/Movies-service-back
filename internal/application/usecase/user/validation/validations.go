package uservalidation

import (
	"net/mail"

	usererror "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil || len(email) > 30 {
		return usererror.ErrUserEmailValidationFailed
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) > 0 && len(password) < 30 {
		return nil
	}
	return usererror.ErrUserPasswordValidationFailed
}

func ValidateUsername(username string) error {
	if len(username) > 0 && len(username) < 30 {
		return nil
	}
	return usererror.ErrUserNameValidationFailed
}

func ValidateUserRegistrationData(data object.UserRegistrationData) error {
	if err := ValidateUsername(data.Username()); err != nil {
		return err
	}

	if err := ValidateEmail(data.Email()); err != nil {
		return err
	}

	if err := ValidatePassword(data.Password()); err != nil {
		return err
	}
	return nil
}

func ValidateAuthenticationData(data object.AuthenticationData) error {
	if err := ValidateEmail(data.Email()); err != nil {
		return err
	}

	if err := ValidatePassword(data.Password()); err != nil {
		return err
	}
	return nil
}
