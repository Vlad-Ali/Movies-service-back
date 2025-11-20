package user

import (
	userdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type UserModel struct {
	ID       string
	Username string
	Email    string
	Password string
}

func (u *UserModel) ToDomain() *userdomain.User {
	user := userdomain.NewUser(u.Username, u.Password, u.Email)
	userID, _ := object.NewUserID(u.ID)
	_ = user.SetID(userID)
	return user
}
