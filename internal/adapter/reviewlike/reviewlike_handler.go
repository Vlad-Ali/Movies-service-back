package reviewlike

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	request2 "github.com/Vlad-Ali/Movies-service-back/internal/adapter/reviewlike/request"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/useridkey"
	error3 "github.com/Vlad-Ali/Movies-service-back/internal/domain/review/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/review/object"
	reviewlikedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike"
	error2 "github.com/Vlad-Ali/Movies-service-back/internal/domain/reviewlike/error"
)

type ReviewLikeHandler struct {
	ReviewLikeService reviewlikedomain.Service
}

func NewReviewLikeHandler(reviewLikeService reviewlikedomain.Service) *ReviewLikeHandler {
	return &ReviewLikeHandler{ReviewLikeService: reviewLikeService}
}

func (rl *ReviewLikeHandler) Like(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewLikeHandler.Like called")

	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("ReviewLikeHandler.Like Error while extracting userID", "Error", err)
		http.Error(w, "Failed to like", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("ReviewLikeHandler.Like Error while reading body", "Error", err)
		http.Error(w, "Failed to like", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var request request2.ReviewLikeRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Error("ReviewLikeHandler.Like Error while unmarshalling body", "Error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	reviewID, err := object.NewReviewID(request.ReviewID)
	if err != nil {
		slog.Error("ReviewLikeHandler.Like Error with reviewID", "Error", err)
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	err = rl.ReviewLikeService.LikeReview(r.Context(), userID, reviewID)

	if err != nil {
		slog.Error("ReviewLikeHandler.Like Error with like review", "Error", err)
		if errors.Is(err, error3.ErrReviewNotFound) {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else if errors.Is(err, error2.ErrReviewLikeAlreadyExists) {
			http.Error(w, "Review like already exists", http.StatusConflict)
		} else {
			http.Error(w, "Review like error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte("successfully liked review"))
	if err != nil {
		slog.Error("ReviewLikeHandler.Like write", "Error", err)
		return
	}
}

func (rl *ReviewLikeHandler) UnLike(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ReviewLikeHandler.UnLike called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("ReviewLikeHandler.UnLike Error while extracting userID", "Error", err)
		http.Error(w, "Failed to unlike", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("ReviewLikeHandler.UnLike Error while reading body", "Error", err)
		http.Error(w, "Failed to like", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var request request2.ReviewLikeRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Error("ReviewLikeHandler.UnLike Error while unmarshalling body", "Error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	reviewID, err := object.NewReviewID(request.ReviewID)
	if err != nil {
		slog.Error("ReviewLikeHandler.UnLike Error with reviewID", "Error", err)
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	err = rl.ReviewLikeService.UnLikeReview(r.Context(), userID, reviewID)
	if err != nil {
		slog.Error("ReviewLikeHandler.UnLike Error with unlike review", "Error", err)
		if errors.Is(err, error3.ErrReviewNotFound) {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else if errors.Is(err, error2.ErrReviewLikeIsNotFound) {
			http.Error(w, "Review like is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Review like error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte("successfully unliked review"))
	if err != nil {
		slog.Error("ReviewLikeHandler.UnLike write", "Error", err)
		return
	}
}
