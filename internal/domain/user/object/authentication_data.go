package object

type AuthenticationData struct {
	password string
	email    string
}

func NewAuthenticationData(password string, email string) AuthenticationData {
	return AuthenticationData{password, email}
}

func (a AuthenticationData) Password() string {
	return a.password
}

func (a AuthenticationData) Email() string {
	return a.email
}
