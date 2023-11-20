package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/reinaldosaraiva/go-api/internal/dto"
	"github.com/reinaldosaraiva/go-api/internal/entity"
	"github.com/reinaldosaraiva/go-api/internal/infra/database"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
	var user dto.LoginDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {	
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {	
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !u.CheckPassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := map[string]interface{}{
		"sub": strconv.FormatUint(uint64(u.ID), 10),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	}
	_, tokenStr, err := jwt.Encode(claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	accessToken := struct{
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenStr,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	
}
