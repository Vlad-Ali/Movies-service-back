package app

import (
	"database/sql"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/jwt"
	movie2 "github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/movie"
	reviewservice "github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/review"
	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/review/modelconfig"
	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/transactionmanager"
	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/user"
	usermovie2 "github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/usermovie"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	userdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
)

type Services struct {
	UserService      userdomain.Service
	MovieService     movie.Service
	UserMovieService usermovie.Service
	TokenService     userdomain.TokenService
	ReviewService    reviewdomain.Service
	ReviewProvider   reviewdomain.Provider
}

func NewServices(db *sql.DB, repos *Repositories, secretKey string, transactionUser transactionmanager.TransactionUser, config modelconfig.ModelConfig) *Services {
	tokenService := jwt.NewJwtService(secretKey)
	userService := user.NewUserService(tokenService, repos.UserRepository, transactionmanager.NewTransactionManager[*userdomain.User](db),
		transactionmanager.NewTransactionManager[*object.AuthResponse](db))
	movieService := movie2.NewMovieService(repos.MovieRepository, transactionmanager.NewTransactionManager[*movie.Movie](db),
		transactionmanager.NewTransactionManager[[]*movie.Movie](db))
	userMovieService := usermovie2.NewUserMovieService(repos.MovieRepository, repos.UserMovieRepository,
		transactionmanager.NewTransactionManager[[]*usermovie.MovieUserInfo](db), transactionmanager.NewTransactionManager[*usermovie.MovieUserInfo](db),
		transactionUser)
	reviewService := reviewservice.NewReviewService(repos.MovieRepository, repos.ReviewRepository, transactionUser, transactionmanager.NewTransactionManager[*reviewdomain.Review](db),
		transactionmanager.NewTransactionManager[[]*reviewdomain.ReviewInfo](db))
	reviewProvider := reviewservice.NewReviewProvider(reviewService, config)
	return &Services{UserService: userService, MovieService: movieService, UserMovieService: userMovieService, TokenService: tokenService, ReviewService: reviewService, ReviewProvider: reviewProvider}
}
