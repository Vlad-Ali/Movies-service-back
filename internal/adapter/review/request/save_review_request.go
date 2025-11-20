package reviewrequest

import (
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type SaveReviewRequest struct {
	Text        string           `json:"text"`
	ReviewYear  int              `json:"review_year"`
	ReviewMonth int              `json:"review_month"`
	ReviewDay   int              `json:"review_day"`
	MovieInfo   object.MovieInfo `json:"movie_info"`
}
