package error

import "errors"

var (
	ErrUserIDCreatingIsNotValid     = errors.New("userID is not valid")
	ErrUserEmailAlreadyExists       = errors.New("user email already exists")
	ErrUserIsNotFound               = errors.New("user not found")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrFailedToRegisterUser         = errors.New("failed to register user")
	ErrFailedToAuthorizeUser        = errors.New("failed to authorize user")
	ErrUserIDAlreadyExists          = errors.New("user id already exists")
	ErrUserNameValidationFailed     = errors.New("user name validation failed")
	ErrUserEmailValidationFailed    = errors.New("user email validation failed")
	ErrUserPasswordValidationFailed = errors.New("user password validation failed")
)
