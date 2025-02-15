package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if err := cfg.db.Reset(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting all users from database", err)
		return
	}
}
