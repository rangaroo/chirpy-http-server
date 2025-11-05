package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rangaroo/chirpy-http-server/internal/auth"
	"github.com/rangaroo/chirpy-http-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
        ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Password  string    `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVals struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't decode parameters", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't hash password", err)
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
