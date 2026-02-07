package main

import (
	"net/http"
	"github.com/Ocidemus/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter,r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
    if err != nil || token == "" {
        respondWithError(w,http.StatusUnauthorized,"cant find token",err)
        return
    }
	_,err = cfg.db.RevokeRefreshToken(r.Context(),token)
	if err != nil{
		respondWithError(w,http.StatusUnauthorized,"invalid token",nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}