package usermovie

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/useridkey"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/usermovie/request"
	response2 "github.com/Vlad-Ali/Movies-service-back/internal/adapter/usermovie/response"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
	error3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie/error"
)

type UserMovieHandler struct {
	userMovieService usermovie.Service
}

func NewUserMovieHandler(userMovieService usermovie.Service) *UserMovieHandler {
	return &UserMovieHandler{userMovieService}
}

func (u *UserMovieHandler) SaveRating(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserMovieHandler.SaveRating called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request: ", "Error", err)
		http.Error(w, "Failed to save rating", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("UserMovieHandler.SaveRating  Error reading body: ", "Error", err)
		http.Error(w, "Failed to save rating", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var saveRatingRequest request.UserMovieSaveRatingRequest
	err = json.Unmarshal(body, &saveRatingRequest)
	if err != nil {
		slog.Error("UserMovieHandler.SaveRating  Error unmarshalling body: ", "Error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = u.userMovieService.SaveRating(r.Context(), userID, saveRatingRequest.MovieInfo, saveRatingRequest.Rating)
	if err != nil {
		slog.Error("UserMovieHandler.SaveRating  Error saving rating: ", "Error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
			return
		} else if errors.Is(err, error3.ErrInvalidRating) {
			http.Error(w, "Invalid rating", http.StatusBadRequest)
			return
		} else if errors.Is(err, error3.ErrUserMovieIsNotFound) {
			http.Error(w, "Movie is not found in this list", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to save rating", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte("Rating saved successfully"))
	if err != nil {
		slog.Error("UserMovieHandler.SaveRating  Error writing body: ", "Error", err)
		return
	}
}

func (u *UserMovieHandler) SaveListType(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserMovieHandler.SaveListType called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request: ", "Error", err)
		http.Error(w, "Failed to save list type", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("UserMovieHandler.SaveListType Error reading body: ", "Error", err)
		http.Error(w, "Failed to save list type", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var saveListTypeRequest request.UserMovieWithListTypeRequest
	err = json.Unmarshal(body, &saveListTypeRequest)
	if err != nil {
		slog.Error("UserMovieHandler.SaveListType Error unmarshalling body: ", "Error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = u.userMovieService.SaveListType(r.Context(), userID, saveListTypeRequest.MovieInfo, saveListTypeRequest.ListType)
	if err != nil {
		slog.Error("UserMovieHandler.SaveListType Error saving list type: ", "Error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
			return
		} else if errors.Is(err, error3.ErrListTypeIsIncorrect) {
			http.Error(w, "Invalid list-type", http.StatusBadRequest)
			return
		} else if errors.Is(err, error3.ErrUserMovieIsNotFound) {
			http.Error(w, "Movie is not found in this list", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to save list-type", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte("List saved successfully"))
	if err != nil {
		slog.Error("UserMovieHandler.SaveListType Error writing body: ", "Error", err)
		return
	}
}

func (u *UserMovieHandler) GetUserMovie(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserMovieHandler.GetUserMovie called")

	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request: ", "Error", err)
		http.Error(w, "Failed to get user movie", http.StatusUnauthorized)
		return
	}

	movieInfo, err := object.GetMovieInfoFromReq(r)

	if err != nil {
		slog.Error("UserMovieHandler.GetUserMovie Error getting parameters: ", "Error", err)
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	listType := r.URL.Query().Get("listType")
	movie, err := u.userMovieService.FindMovieByUser(r.Context(), userID, movieInfo, listType)
	if err != nil {
		slog.Error("UserMovieHandler.GetUserMovie  Error finding movie: ", "Error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else if errors.Is(err, error3.ErrUserMovieIsNotFound) {
			http.Error(w, "This movie is not found in this list", http.StatusBadRequest)
		} else if errors.Is(err, error3.ErrListTypeIsIncorrect) {
			http.Error(w, "Invalid list-type", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to get user movie", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		slog.Error("UserMovieHandler.GetUserMovie  Error writing body: ", "Error", err)
		return
	}
}

func (u *UserMovieHandler) GetUserMovies(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserMovieHandler.GetUserMovies called")

	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request: ", "Error", err)
		http.Error(w, "Failed to get user movie", http.StatusUnauthorized)
		return
	}

	listType := r.URL.Query().Get("listType")
	slog.Debug("Got list type", "listType", listType)
	movies, err := u.userMovieService.FindMoviesByUserAndListType(r.Context(), userID, listType)
	if err != nil {
		slog.Error("UserMovieHandler.GetUserMovies  Error finding movies: ", "Error", err)
		if errors.Is(err, error3.ErrListTypeIsIncorrect) {
			http.Error(w, "Invalid list-type", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to get user movies", http.StatusInternalServerError)
		}
		return
	}

	response := response2.UserMoviesResponse{UserMovies: movies}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("UserMovieHandler.GetUserMovies Error writing body: ", "Error", err)
		return
	}
}
