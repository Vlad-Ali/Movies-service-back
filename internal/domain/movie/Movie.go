package movie

import (
	"strings"
	"time"

	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type Movie struct {
	id          object.MovieID
	Title       string
	Description string
	ReleaseDate time.Time
	Director    string
	Actors      []string
	Genres      []string
	Rating      float64
}

func NewMovie(title, description string, releaseDate time.Time, director string, actors, genres []string, rating float64) *Movie {
	return &Movie{
		id:          object.MovieID{},
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		ReleaseDate: releaseDate,
		Director:    strings.TrimSpace(director),
		Actors:      actors,
		Genres:      genres,
		Rating:      rating,
	}
}

func (m *Movie) ID() object.MovieID {
	return m.id
}

func (m *Movie) SetID(id object.MovieID) error {
	if m.id.IsEmpty() {
		m.id = id
		return nil
	}
	return error2.ErrMovieIDAlreadyExists
}
