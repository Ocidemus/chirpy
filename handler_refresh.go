package main

import (
	"net/http"
	"time"

	"github.com/Ocidemus/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnval struct{
		Token string `json:"token"`
	}
    token, err := auth.GetBearerToken(r.Header)
    if err != nil || token == "" {
        respondWithError(w,http.StatusUnauthorized,"cant find token",err)
        return
    }
	user,err := cfg.db.GetUserFromRefreshToken(r.Context(),token)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"invalid token",nil)
		return
	}
	new_token, err := auth.MakeJWT(user.ID, cfg.secret ,time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
}
		
    respondWithJSON(w, http.StatusOK,returnval{
		Token: new_token,
	})
}