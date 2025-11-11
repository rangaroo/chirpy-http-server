package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rangaroo/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could't parse the chirpID", err)
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could't parse the token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could't get chirp", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "You can not delete other's chirps", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
