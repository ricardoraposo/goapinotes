package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/ricardoraposo/api-again/internal/database"
	"github.com/ricardoraposo/api-again/internal/dto"
	"github.com/ricardoraposo/api-again/internal/entity"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{UserDB: db, Jwt: jwt, JwtExpiresIn: jwtExpiresIn}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var u dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UserDB.FindByEmail(u.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !user.ValidatePassword(u.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.Jwt.Encode(map[string]interface{}{
		"id":        user.ID,
		"expiresIn": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"access_token": tokenString})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u dto.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(u.Name, u.Email, u.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
