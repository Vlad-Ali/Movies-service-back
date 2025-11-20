package request

import object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"

type UserMovieWithListTypeRequest struct {
	MovieInfo object2.MovieInfo `json:"movie_info"`
	ListType  string            `json:"list_type"`
}
