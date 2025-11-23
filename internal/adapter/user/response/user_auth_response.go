package response

type UserAuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}
