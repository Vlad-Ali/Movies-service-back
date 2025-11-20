package review

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	reviewrequest "github.com/Vlad-Ali/Movies-service-back/internal/adapter/review/request"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/review/response"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/useridkey"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie/object"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	error3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/error"
)

type ReviewHandler struct {
	reviewService  reviewdomain.Service
	reviewProvider reviewdomain.Provider
}

func NewReviewHandler(reviewService reviewdomain.Service, provider reviewdomain.Provider) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService, reviewProvider: provider}
}

func (rh *ReviewHandler) SaveReview(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewHandler.SaveReview called")

	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request", "error", err)
		http.Error(w, "Failed to save review", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error while reading body", "error", err)
		http.Error(w, "Failed to save review", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var saveRequest reviewrequest.SaveReviewRequest
	err = json.Unmarshal(body, &saveRequest)
	if err != nil {
		slog.Error("Error while unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	date := time.Date(saveRequest.ReviewYear, time.Month(saveRequest.ReviewMonth), saveRequest.ReviewDay, 0, 0, 0, 0, time.UTC)
	err = rh.reviewService.SaveReview(r.Context(), userID, saveRequest.MovieInfo, saveRequest.Text, date)
	if err != nil {
		slog.Error("Error while saving review", "error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else if errors.Is(err, error3.ErrReviewTextValidationError) {
			http.Error(w, "Text validation error", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to save review", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte("Successfully saved review"))
	if err != nil {
		slog.Error("Error while writing body", "error", err)
		return
	}
}

func (rh *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewHandler.DeleteReview called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request", "error", err)
		http.Error(w, "Failed to delete review", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error while reading body", "error", err)
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var movieInfo object.MovieInfo
	err = json.Unmarshal(body, &movieInfo)
	if err != nil {
		slog.Error("Error while unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = rh.reviewService.DeleteReview(r.Context(), userID, movieInfo)
	if err != nil {
		slog.Error("Error while deleting review", "error", err)
		if errors.Is(err, error3.ErrReviewNotFound) {
			http.Error(w, "Review is not found", http.StatusNotFound)
		} else if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte("Successfully deleted review"))
	if err != nil {
		slog.Error("Error while writing body", "error", err)
		return
	}
}

func (rh *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewHandler.GetReview called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error while extracting user id from request", "error", err)
		http.Error(w, "Failed to get review", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error while reading body", "error", err)
		http.Error(w, "Failed to get review", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var movieInfo object.MovieInfo
	err = json.Unmarshal(body, &movieInfo)
	if err != nil {
		slog.Error("Error while unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	review, err := rh.reviewService.GetUserReview(r.Context(), userID, movieInfo)
	if err != nil {
		slog.Error("Error while getting review", "error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get review", http.StatusInternalServerError)
		}
		return
	}

	reviewResponse := response.GetReviewResponse{Text: review.Text(), ReviewYear: review.WritingDate().Year(), ReviewMonth: int(review.WritingDate().Month()), ReviewDay: review.WritingDate().Day()}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviewResponse)
	if err != nil {
		slog.Error("Error while writing body", "error", err)
		return
	}
}

func (rh *ReviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewHandler.GetReviews called")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error while reading body", "error", err)
		http.Error(w, "Failed to get reviews", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var movieInfo object.MovieInfo
	err = json.Unmarshal(body, &movieInfo)
	if err != nil {
		slog.Error("Error while unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	reviews, err := rh.reviewService.GetUserReviews(r.Context(), movieInfo)
	if err != nil {
		slog.Error("Error while getting reviews", "error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get reviews", http.StatusInternalServerError)
		}
		return
	}

	getResponse := response.GetReviewsResponse{Reviews: reviews}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(getResponse)
	if err != nil {
		slog.Error("Error while writing body", "error", err)
		return
	}
}

func (rh *ReviewHandler) GetSummaryReviews(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewHandler.GetSummaryReviews called")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error while reading body", "error", err)
		http.Error(w, "Failed to get reviews", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var movieInfo object.MovieInfo
	err = json.Unmarshal(body, &movieInfo)
	if err != nil {
		slog.Error("Error while unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	summary, err := rh.reviewProvider.ProvideMovieReviews(r.Context(), movieInfo)
	if err != nil {
		slog.Error("Error while getting summary", "error", err)
		if errors.Is(err, error2.ErrMovieIsNotFound) {
			http.Error(w, "Movie is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		}
		return
	}

	summaryResponse := response.GetSummaryResponse{Summary: summary}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(summaryResponse)
	if err != nil {
		slog.Error("Error while writing body", "error", err)
		return
	}
	slog.Info("Successfully got summary")
}
