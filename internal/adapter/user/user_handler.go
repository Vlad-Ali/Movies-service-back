package user

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/user/request"
	userresponse "github.com/Vlad-Ali/Movies-service-back/internal/adapter/user/response"
	"github.com/Vlad-Ali/Movies-service-back/internal/adapter/useridkey"
	userdomain "github.com/Vlad-Ali/Movies-service-back/internal/domain/user"
	usererror "github.com/Vlad-Ali/Movies-service-back/internal/domain/user/error"
	"github.com/Vlad-Ali/Movies-service-back/internal/domain/user/object"
)

type UserHandler struct {
	userService userdomain.Service
}

func NewUserHandler(userService userdomain.Service) *UserHandler {
	return &UserHandler{userService}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserHandler.Register called")
	var registerRequest request.UserRegisterRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error reading body", "error", err)
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		slog.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	registerData := object.NewUserRegistrationData(registerRequest.Username, registerRequest.Password, registerRequest.Email)
	user, err := u.userService.Register(r.Context(), registerData)
	if err != nil {
		slog.Error("Error registering user", "error", err)
		if errors.Is(err, usererror.ErrUserEmailAlreadyExists) {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else if errors.Is(err, usererror.ErrUserNameValidationFailed) {
			http.Error(w, "Username is empty", http.StatusBadRequest)
		} else if errors.Is(err, usererror.ErrUserPasswordValidationFailed) {
			http.Error(w, "Password is empty", http.StatusBadRequest)
		} else if errors.Is(err, usererror.ErrUserEmailValidationFailed) {
			http.Error(w, "Email is invalid", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to register", http.StatusInternalServerError)
		}
		return
	}

	response := userresponse.UserRegisterResponse{Username: user.Username(), Email: user.Email()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}

func (u *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserHandler.Authenticate called")
	var authRequest request.UserAuthRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Error reading body", "error", err)
		http.Error(w, "Failed to Authenticate", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &authRequest)
	if err != nil {
		slog.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	authData := object.NewAuthenticationData(authRequest.Password, authRequest.Email)
	token, err := u.userService.Authenticate(r.Context(), authData)
	if err != nil {
		slog.Error("Error authenticating user", "error", err)
		if errors.Is(err, usererror.ErrUserPasswordValidationFailed) || errors.Is(err, usererror.ErrUserEmailValidationFailed) {
			http.Error(w, "Invalid input", http.StatusBadRequest)
		} else if errors.Is(err, usererror.ErrUserIsNotFound) || errors.Is(err, usererror.ErrInvalidPassword) {
			http.Error(w, "Email or password are invalid", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		}
		return
	}

	response := userresponse.UserAuthResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("UserHandler.GetUser called")
	userID, err := useridkey.ExtractUserIdFromReq(r)
	if err != nil {
		slog.Error("Error extracting user id", "error", err)
		http.Error(w, "Failed to get user", http.StatusUnauthorized)
		return
	}

	user, err := u.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		slog.Error("Error getting user", "error", err)
		if errors.Is(err, usererror.ErrUserIsNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
		}
		return
	}

	response := userresponse.UserGetResponse{Username: user.Username(), Email: user.Email()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		return
	}
}
