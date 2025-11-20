package object

type UserRegistrationData struct {
	username string
	password string
	email    string
}

func NewUserRegistrationData(name string, password string, email string) UserRegistrationData {
	return UserRegistrationData{username: name, password: password, email: email}
}

func (u UserRegistrationData) Username() string {
	return u.username
}

func (u UserRegistrationData) Password() string {
	return u.password
}

func (u UserRegistrationData) Email() string {
	return u.email
}
