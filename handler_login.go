package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ocidemus/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handle_login(w http.ResponseWriter,r *http.Request){
	type body struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type returnval struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	user,err:=cfg.db.GetUserByEmail(r.Context(),params.Email)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"Incorrect email or password",nil)
		return
	}
	passcheck,err := auth.CheckPasswordHash(params.Password,user.HashedPassword)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"Incorrect email or password",nil)
		return
	}
	if passcheck == false {
		respondWithError(w,http.StatusUnauthorized,"Incorrect email or password",nil)
		return
	}
	respondWithJSON(w, http.StatusOK, returnval{
		ID:user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
	
}