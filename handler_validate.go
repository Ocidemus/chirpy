package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ocidemus/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		User_id string `json:"user_id"`
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

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleaned := cleanprofanity(params.Body)
	parsedID, err := uuid.Parse(params.User_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user_id", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: uuid.NullUUID{
			UUID:  parsedID,
			Valid: true,
		},
})
	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: cleaned,
		UserID: parsedID,
	})
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
