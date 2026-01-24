package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanprofanity(params.Body),
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
