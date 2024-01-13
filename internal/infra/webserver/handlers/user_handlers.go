package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/pedrohenmonteiro/golang-api/internal/dto"
	"github.com/pedrohenmonteiro/golang-api/internal/entity"
	"github.com/pedrohenmonteiro/golang-api/internal/infra/database"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJWTInput

	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

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

	if !u.ComparePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claimsMap := map[string]interface{}{
		"sub":   u.ID.String(),
		"exp":   time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
		"email": u.Email,
	}

	_, tokenString, _ := jwt.Encode(claimsMap)

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)

}
