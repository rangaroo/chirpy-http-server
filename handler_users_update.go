package main

import (
	"encoding/json"
	"net/http"

	"github.com/rangaroo/chirpy-http-server/internal/auth"
	"github.com/rangaroo/chirpy-http-server/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVals struct {
		User
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could't parse the token", err)
		return
	}

	UserID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't decode parameters", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't hash password", err)
	}

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
		ID:             UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
