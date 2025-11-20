package jwt

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	jwtclaims "github.com/Vlad-Ali/Movies-service-back/internal/application/dto/jwt"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	usererror "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretKey string
}

func NewJwtService(secretKey string) *JwtService {
	return &JwtService{secretKey: secretKey}
}

func (j *JwtService) GenerateToken(ctx context.Context, user *user.User) (string, error) {
	now := time.Now()
	duration := 24 * time.Hour
	claims := jwtclaims.JWTClaims{
		Username: user.Username(),
		Email:    user.Email(),
		UserID:   user.ID().ID(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	slog.Debug("generate token with id", "userID", user.ID().ID())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JwtService) ValidateToken(ctx context.Context, token string) (object.UserID, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtclaims.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		slog.Error("validate token error", "error", err)
		return object.UserID{}, usererror.ErrFailedToAuthorizeUser
	}

	if claims, ok := jwtToken.Claims.(*jwtclaims.JWTClaims); ok && jwtToken.Valid {
		userID, err := object.NewUserID(claims.UserID)
		if err != nil {
			slog.Error("userID is incorrect", "error", err)
			return object.UserID{}, err
		}
		slog.Debug("validation of token is successful with ID", "ID", claims.UserID)
		return userID, nil
	}

	slog.Error("validate token error", "error", err)
	return object.UserID{}, usererror.ErrFailedToAuthorizeUser
}
