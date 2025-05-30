package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respBody := errorResponse{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	// check profanity
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(params.Body, " ")
	var cleanedWords []string
	for _, word := range words {
		if slices.Contains(profanity, strings.ToLower(word)) {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}

	respBody := successResponse{
		CleanedBody: strings.Join(cleanedWords, " "),
	}
	respondWithJSON(w, 200, respBody)
}
