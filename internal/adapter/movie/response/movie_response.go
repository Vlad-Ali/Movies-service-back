package movieresponse

import "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"

type MovieResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Month       int      `json:"month"`
	Day         int      `json:"day"`
	Director    string   `json:"director"`
	Actors      []string `json:"actors"`
	Genres      []string `json:"genres"`
	Rating      float64  `json:"rating"`
}

func NewMovieResponse(movie *movie.Movie) MovieResponse {
	return MovieResponse{
		Title:       movie.Title,
		Description: movie.Description,
		Year:        movie.ReleaseDate.Year(),
		Month:       int(movie.ReleaseDate.Month()),
		Day:         movie.ReleaseDate.Day(),
		Director:    movie.Director,
		Actors:      movie.Actors,
		Genres:      movie.Genres,
		Rating:      movie.Rating,
	}
}
