package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ocidemus/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polka_hook(w http.ResponseWriter,r *http.Request){
	type body struct {
		Event string `json:"event"`
		Data struct {
			User_id uuid.UUID `json:"user_id"`
		}
	}
	api_key,err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w,http.StatusUnauthorized,"key not found",err)
		return
	}
	if api_key != cfg.polka_key {
		respondWithError(w,http.StatusUnauthorized,"unauthorized",err)
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_,err = cfg.db.UpgradeToChirpyRed(r.Context(),params.Data.User_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	
}