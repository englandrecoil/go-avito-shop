package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/englandrecoil/go-avito-shop/internal/auth"
	"github.com/englandrecoil/go-avito-shop/internal/database"
	"github.com/lib/pq"
)

func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {
	type registerRequestParams struct {
		Username string `json:"username"`
		Passowrd string `json:"password"`
	}
	type registerResponseParams struct {
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Balance   int       `json:"balance"`
	}
	reqUser := registerRequestParams{}

	// get user's request params
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqUser); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(reqUser.Passowrd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// store data in db
	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username:       reqUser.Username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				respondWithError(w, http.StatusBadRequest, "This username's already in use", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
			return
		}
	}

	// send response
	respondWithJSON(w, http.StatusCreated, registerResponseParams{
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Balance:   int(dbUser.Balance),
	})
}
