package usermovie

import (
	"time"
)

type MovieUserInfo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Director    string    `json:"director"`
	Actors      []string  `json:"actors"`
	Genres      []string  `json:"genres"`
	Rating      float64   `json:"rating"`

	ListType   ListType `json:"list_type"`
	UserRating int      `json:"user_rating"`
}
