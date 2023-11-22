package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/reinaldosaraiva/go-api/internal/dto"
	"github.com/reinaldosaraiva/go-api/internal/entity"
	"github.com/reinaldosaraiva/go-api/internal/infra/database"
)

type Error struct{
	Message string `json:"message"`
}
type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}
// Get a user by JWT godoc
// @Summary Get a user by JWT user
// @Description Get a user by JWT user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.LoginDTO true "user request"
// @Success 200
// @Failure 500 {object} Error
// @Failure 404 {object} Error
// @Router /users/generate_token [post]
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
		w.WriteHeader(http.StatusNotFound)
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
	accessToken := dto.GetJWTOutput{AccessToken: tokenStr}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserDTO true "user request"
// @Success 201
// @Failure 400 {object} Error
// @Router /users [post]
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

//Get User with email by body request
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.UserDB.FindByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}
