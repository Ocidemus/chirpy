package main

import "net/http"



func (cfg *apiConfig) reset(w http.ResponseWriter,r *http.Request){
	if cfg.platform != "dev"{
		respondWithError(w,http.StatusForbidden,"FORBIDDEN",nil)
		return
	}
	err := cfg.db.Reset(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset user", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK) // 200
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
