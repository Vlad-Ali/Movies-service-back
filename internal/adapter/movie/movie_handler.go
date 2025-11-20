package movie

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	movieresponse "github.com/Vlad-Ali/Movies-service-back/internal/adapter/movie/response"
	moviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
)

type MovieHandler struct {
	movieService moviedomain.Service
}

func NewMovieHandler(movieService moviedomain.Service) *MovieHandler {
	return &MovieHandler{movieService}
}

func (m *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	slog.Debug("MovieHandler.GetMovie called")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("MovieHandler.GetMovie error reading Body", slog.String("err", err.Error()))
		http.Error(w, "Failed to get movie", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var movieInfo object.MovieInfo
	err = json.Unmarshal(body, &movieInfo)
	if err != nil {
		slog.Error("MovieHandler.GetMovie error unmarshalling body", slog.String("err", err.Error()))
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
		return
	}

	movie, err := m.movieService.FindByReleaseDateAndTitle(r.Context(), movieInfo)
	if err != nil {
		slog.Error("MovieHandler.GetMovie Error finding movie", slog.String("err", err.Error()))
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get movie", http.StatusInternalServerError)
		}
		return
	}
	movieResponse := movieresponse.NewMovieResponse(movie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(movieResponse); err != nil {
		slog.Error("MovieHandler.GetMovie error encoding response", slog.String("err", err.Error()))
		return
	}
}

func (m *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	slog.Debug("MovieHandler.GetMovies called")

	movies, err := m.movieService.GetAll(r.Context())
	if err != nil {
		slog.Error("MovieHandler.GetMovies Error finding movies", slog.String("err", err.Error()))
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	moviesResponse := movieresponse.MoviesResponse{Movies: make([]movieresponse.MovieResponse, 0)}
	for _, movie := range movies {
		movieResponse := movieresponse.NewMovieResponse(movie)
		moviesResponse.Movies = append(moviesResponse.Movies, movieResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(moviesResponse); err != nil {
		slog.Error("MovieHandler.GetMovies Error encoding response", slog.String("err", err.Error()))
		return
	}
}
