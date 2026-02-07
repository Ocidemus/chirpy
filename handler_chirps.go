package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ocidemus/chirpy/internal/auth"
	"github.com/Ocidemus/chirpy/internal/database"
	"github.com/google/uuid"
)
func (cfg *apiConfig) handlechirp (w http.ResponseWriter, r *http.Request){
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirpID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_id", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(),parsedID)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "Couldn't find any chirps", err)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID.UUID,
	})
}


func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed auth header", err)
		return
	}
	valid_id,err:= auth.ValidateJWT(token,cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid id found", err)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleaned := cleanprofanity(params.Body)
	// parsedID, err := uuid.Parse(params.User_id)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "invalid user_id", err)
	// 	return
	// }
    // if parsedID != valid_id {
	// 	respondWithError(w, http.StatusUnauthorized, "no valid id found", err)
	// 	return
	// }
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: uuid.NullUUID{
			UUID:  valid_id,
			Valid: true,
		},
})
	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: cleaned,
		UserID: chirp.UserID.UUID,
	})
}

func (cfg *apiConfig) reqchirp(w http.ResponseWriter,r *http.Request){
	type result struct{
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't find any chirps", err)
		return
	}
	arr := []result{}
	for _,chirp := range dbChirps{
		arr = append(arr, result{
			ID:chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID.UUID,
		})
	}
	respondWithJSON(w, http.StatusOK, arr)

}

func cleanprofanity(text string) string {
	words := strings.Fields(text)
	for index,word := range words {
		if strings.ToLower(word) == "kerfuffle"{
			words[index]="****"
	}
	if strings.ToLower(word) == "sharbert"{
		words[index]="****"
	}

	if strings.ToLower(word) == "fornax"{
		words[index]="****"
	}
}
	return strings.Join(words, " ")
}
