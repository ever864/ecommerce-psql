package user

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ever864/ecommerce-psql/service/auth"
	"github.com/ever864/ecommerce-psql/types"
	"github.com/ever864/ecommerce-psql/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if the user exists
	existingUser, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error retrieving user: %v", err))
		return
	}

	if existingUser == nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user with email %s does not exist", payload.Email))
		return
	}

	// check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(payload.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid password"))
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	godenv := os.Getenv("JWT_SECRET")

	token, err := auth.CreateJWT([]byte(godenv), existingUser.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// ckeck if the user exists
	existingUser, err := h.store.GetUserByEmail(payload.Email)
	if err == nil && existingUser != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create the user
	user := types.User{
		Email:    payload.Email,
		Password: hashedPassword,
	}

	if err := h.store.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// if it doent we create the new user
	log.Println(err)

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}
