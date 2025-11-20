package object

import (
	"github.com/google/uuid"

	usererror "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/error"
)

type UserID struct {
	id string
}

func NewUserID(s string) (UserID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return UserID{}, usererror.ErrUserIDCreatingIsNotValid
	}
	return UserID{id: s}, nil
}

func (u UserID) ID() string {
	return u.id
}

func (u UserID) IsEmpty() bool {
	return u.id == ""
}
