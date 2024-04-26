package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func profanityFilter(text string) string {
	badWords := [3]string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(text)
	for i, word := range words {
		wordLower := strings.ToLower(word)
		for _, badWord := range badWords {
			if wordLower == badWord {
				words[i] = "****"
			}
		}
	}
	cleanText := strings.Join(words, " ")
	return cleanText
}

func jsonError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
