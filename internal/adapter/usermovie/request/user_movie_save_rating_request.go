package request

import object2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"

type UserMovieSaveRatingRequest struct {
	MovieInfo object2.MovieInfo `json:"movie_info"`
	Rating    int               `json:"rating"`
}
