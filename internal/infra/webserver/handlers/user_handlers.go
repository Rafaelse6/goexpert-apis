package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Rafaelse6/goexpert/9-APIS/internal/dto"
	"github.com/Rafaelse6/goexpert/9-APIS/internal/entity"
	"github.com/Rafaelse6/goexpert/9-APIS/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwtValue := r.Context().Value("jwt")
	jwtExpiresValue := r.Context().Value("JwtExpiresIn")

	if jwtValue == nil || jwtExpiresValue == nil {
		http.Error(w, "JWT authentication not configured", http.StatusInternalServerError)
		return
	}

	jwt, ok := jwtValue.(*jwtauth.JWTAuth)
	if !ok {
		log.Println("Erro: JWTAuth inválido")
		http.Error(w, "Invalid JWT authentication", http.StatusInternalServerError)
		return
	}

	jwtExpiresIn, ok := jwtExpiresValue.(int)
	if !ok {
		log.Println("Erro: JwtExpiresIn inválido")
		http.Error(w, "Invalid JWT expiration format", http.StatusInternalServerError)
		return
	}

	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		log.Println("Usuário não encontrado:", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid email or password"})
		return
	}

	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid email or password"})
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user data"})
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
