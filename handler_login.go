package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rangaroo/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type returnVals struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	secs := params.ExpiresInSeconds
	expiresIn := time.Duration(secs) * time.Second
	if expiresIn > time.Hour || secs <= 0 {
		expiresIn = time.Hour
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't create a token string", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: tokenString,
	})
}
