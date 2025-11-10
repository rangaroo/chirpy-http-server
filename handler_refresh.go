package main

import (
	"net/http"
	"time"

	"github.com/rangaroo/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could't parse the token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could't get the user for refresh token", err)
		return
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: tokenString,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could't parse the token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't revoke the token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
