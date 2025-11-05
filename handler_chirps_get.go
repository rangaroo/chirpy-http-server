package main

import (
	"net/http"
	"github.com/google/uuid"
	"database/sql"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("chirpID")

	if len(path) == 0 {
		chirps, err := cfg.db.GetChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could't get chirps", err)
			return
		}

		responseChirps := []Chirp{}
		for _, chirp := range chirps {
			responseChirps = append(responseChirps, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}

		respondWithJSON(w, http.StatusOK, responseChirps)
	} else {
		chirpID, err := uuid.Parse(path)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could't parse the chirpID", err)
		}

		chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Could't find chirp", err)
			return
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could't get chirp", err)
			return
		}

		respondWithJSON(w, http.StatusOK, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
}
