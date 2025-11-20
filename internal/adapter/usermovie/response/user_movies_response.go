package response

import "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"

type UserMoviesResponse struct {
	UserMovies []*usermovie.MovieUserInfo `json:"userMovies"`
}
