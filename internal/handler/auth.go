package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/p2p-messanger/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) AuthRouters(r chi.Router) {
	r.Post("/register", ah.Register)
	r.Post("/login", ah.Login)
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
}

type RegisterRequest struct {
	Login    string `json:"user"`
	Password string `json:"password"`
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalit Request body", http.StatusBadRequest)
		return
	}

	user, err := ah.authService.RegistUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token, err := ah.authService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID    int64  `json:"id"`
			Login string `json:"login"`
		}{
			ID:    user.ID,
			Login: user.Login,
		},
	}

	w.Header().Set("Conntent-type", "aplication/json")
	json.NewEncoder(w).Encode(response)
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Bydy", http.StatusBadRequest)
		return
	}

	user, err := ah.authService.LoginUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	token, err := ah.authService.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID    int64  `json:"id"`
			Login string `json:"login"`
		}{
			ID:    user.ID,
			Login: user.Login,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
