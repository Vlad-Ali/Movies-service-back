package app

import (
	"net/http"

	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/middleware"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/movie"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/review"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/user"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/usermovie"
)

type Handlers struct {
	UserHandler      *user.UserHandler
	MovieHandler     *movie.MovieHandler
	UserMovieHandler *usermovie.UserMovieHandler
	AuthHandler      *middleware.AuthMiddleware
	ReviewHandler    *review.ReviewHandler
}

func NewHandlers(services *Services) *Handlers {
	userHandler := user.NewUserHandler(services.UserService)
	movieHandler := movie.NewMovieHandler(services.MovieService)
	userMovieHandler := usermovie.NewUserMovieHandler(services.UserMovieService)
	tokenHandler := middleware.NewAuthMiddleware(services.TokenService)
	reviewHandler := review.NewReviewHandler(services.ReviewService, services.ReviewProvider)
	return &Handlers{UserHandler: userHandler, MovieHandler: movieHandler, UserMovieHandler: userMovieHandler, AuthHandler: tokenHandler,
		ReviewHandler: reviewHandler}
}

func (h *Handlers) registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/user/register", h.UserHandler.Register)
	mux.HandleFunc("POST /api/user/auth", h.UserHandler.Authenticate)
	mux.HandleFunc("GET /api/user", h.UserHandler.GetUser)

	mux.HandleFunc("GET /api/movie", h.MovieHandler.GetMovie)
	mux.HandleFunc("GET /api/movie/all", h.MovieHandler.GetMovies)

	mux.HandleFunc("PATCH /api/user/movie/rating", h.UserMovieHandler.SaveRating)
	mux.HandleFunc("PATCH /api/user/movie/list", h.UserMovieHandler.SaveListType)
	mux.HandleFunc("GET /api/user/movie", h.UserMovieHandler.GetUserMovie)
	mux.HandleFunc("GET /api/user/movie/all", h.UserMovieHandler.GetUserMovies)

	mux.HandleFunc("PUT /api/user/movie/review", h.ReviewHandler.SaveReview)
	mux.HandleFunc("DELETE /api/user/movie/review", h.ReviewHandler.DeleteReview)
	mux.HandleFunc("GET /api/user/movie/review", h.ReviewHandler.GetReview)
	mux.HandleFunc("GET /api/movie/review/all", h.ReviewHandler.GetReviews)
	mux.HandleFunc("GET /api/movie/summary", h.ReviewHandler.GetSummaryReviews)

	mainHandler := h.AuthHandler.Authorize(mux)

	return mainHandler
}
