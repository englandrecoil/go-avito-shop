package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/englandrecoil/go-avito-shop/internal/auth"
)

type authResponseParams struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerAuth(w http.ResponseWriter, r *http.Request) {

	authUser := CredentialsRequestParams{}

	// get user's request params
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&authUser); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	// validate request's params
	if authUser.Username == "" || authUser.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Username and password are required", errors.New("missing parameters"))
	}

	// get user by username from db
	dbUser, err := cfg.db.GetUserByUsername(r.Context(), authUser.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect username or password", err)
		return
	}

	// authentication
	if err = auth.CheckPasswordHash(authUser.Password, dbUser.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password", err)
		return
	}

	// make JWT
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT for user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, authResponseParams{
		Token: accessToken,
	})
}
