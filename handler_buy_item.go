package main

import (
	"net/http"

	"github.com/englandrecoil/go-avito-shop/internal/auth"
	"github.com/englandrecoil/go-avito-shop/internal/database"
)

func (cfg *apiConfig) handlerBuyItem(w http.ResponseWriter, r *http.Request) {
	// get token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find bearer token", err)
		return
	}

	// validate token
	jwtID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate bearer token", err)
		return
	}

	itemNameParam := r.PathValue("item")
	if itemNameParam == "" {
		respondWithError(w, http.StatusBadRequest, "No item specified", err)
		return
	}

	// get item data from db
	item, err := cfg.db.GetItemByName(r.Context(), itemNameParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find item", err)
		return
	}

	// process transaction
	_, err = cfg.db.PurchaseItemByID(r.Context(), database.PurchaseItemByIDParams{
		UserID: jwtID,
		ItemID: item.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't buy item", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
