package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ocidemus/chirpy/internal/auth"
	"github.com/Ocidemus/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) update_email(w http.ResponseWriter,r *http.Request){
	type body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type returnval struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	token, err := auth.GetBearerToken(r.Header)
    if err != nil || token == "" {
        respondWithError(w,http.StatusUnauthorized,"cant find token",err)
        return
    }
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w,http.StatusUnauthorized,"cant find token",err)
        return
	}
	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	dbuser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})	
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	respondWithJSON(w,http.StatusOK,returnval{
		ID :dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Email: dbuser.Email,
	})

}
