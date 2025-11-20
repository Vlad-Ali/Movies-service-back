package app

import (
	"database/sql"

	moviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/movie"
	reviewdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/review"
	userdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	usermoviedomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/usermovie"
	"github.com/Vlad-Ali/Movies-service-back/internal/infrastruture/movie"
	reviewrepo "github.com/Vlad-Ali/Movies-service-back/internal/infrastruture/review"
	"github.com/Vlad-Ali/Movies-service-back/internal/infrastruture/user"
	"github.com/Vlad-Ali/Movies-service-back/internal/infrastruture/usermovie"
)

type Repositories struct {
	MovieRepository     moviedomain.Repository
	UserRepository      userdomain.Repository
	UserMovieRepository usermoviedomain.Repository
	ReviewRepository    reviewdomain.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{MovieRepository: movie.NewMovieRepository(db), UserRepository: user.NewUserRepository(db), UserMovieRepository: usermovie.NewUserMovieRepository(db),
		ReviewRepository: reviewrepo.NewReviewRepository(db)}
}
